// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package import_process

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/mattermost/mattermost-server/v5/app"
	"github.com/mattermost/mattermost-server/v5/app/request"
	"github.com/mattermost/mattermost-server/v5/jobs"
	tjobs "github.com/mattermost/mattermost-server/v5/jobs/interfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/shared/mlog"
)

func init() {
	app.RegisterJobsImportProcessInterface(func(s *app.Server) tjobs.ImportProcessInterface {
		a := app.New(app.ServerConnector(s))
		return &ImportProcessInterfaceImpl{a}
	})
}

type ImportProcessInterfaceImpl struct {
	app *app.App
}

type ImportProcessWorker struct {
	name        string
	stopChan    chan struct{}
	stoppedChan chan struct{}
	jobsChan    chan model.Job
	jobServer   *jobs.JobServer
	app         *app.App
	appContext  *request.Context
}

func (i *ImportProcessInterfaceImpl) MakeWorker() model.Worker {
	return &ImportProcessWorker{
		name:        "ImportProcess",
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
		jobsChan:    make(chan model.Job),
		jobServer:   i.app.Srv().Jobs,
		app:         i.app,
		appContext:  &request.Context{},
	}
}

func (w *ImportProcessWorker) JobChannel() chan<- model.Job {
	return w.jobsChan
}

func (w *ImportProcessWorker) Run() {
	mlog.Debug("Worker started", mlog.String("worker", w.name))

	defer func() {
		mlog.Debug("Worker finished", mlog.String("worker", w.name))
		close(w.stoppedChan)
	}()

	for {
		select {
		case <-w.stopChan:
			mlog.Debug("Worker received stop signal", mlog.String("worker", w.name))
			return
		case job := <-w.jobsChan:
			mlog.Debug("Worker received a new candidate job.", mlog.String("worker", w.name))
			w.doJob(&job)
		}
	}
}

func (w *ImportProcessWorker) Stop() {
	mlog.Debug("Worker stopping", mlog.String("worker", w.name))
	close(w.stopChan)
	<-w.stoppedChan
}

func (w *ImportProcessWorker) doJob(job *model.Job) {
	if claimed, err := w.jobServer.ClaimJob(job); err != nil {
		mlog.Warn("Worker experienced an error while trying to claim job",
			mlog.String("worker", w.name),
			mlog.String("job_id", job.Id),
			mlog.String("error", err.Error()))
		return
	} else if !claimed {
		return
	}

	importFileName, ok := job.Data["import_file"]
	if !ok {
		appError := model.NewAppError("ImportProcessWorker", "import_process.worker.do_job.missing_file", nil, "", http.StatusBadRequest)
		w.setJobError(job, appError)
		return
	}

	importFilePath := filepath.Join(*w.app.Config().ImportSettings.Directory, importFileName)
	if ok, err := w.app.FileExists(importFilePath); err != nil {
		w.setJobError(job, err)
		return
	} else if !ok {
		appError := model.NewAppError("ImportProcessWorker", "import_process.worker.do_job.file_exists", nil, "", http.StatusBadRequest)
		w.setJobError(job, appError)
		return
	}

	importFileSize, appErr := w.app.FileSize(importFilePath)
	if appErr != nil {
		w.setJobError(job, appErr)
		return
	}

	importFile, appErr := w.app.FileReader(importFilePath)
	if appErr != nil {
		w.setJobError(job, appErr)
		return
	}
	defer importFile.Close()

	// extract the contents of the zipped file.
	zipReader, err := zip.NewReader(importFile.(io.ReaderAt), importFileSize)
	if err != nil {
		panic(err)
	}

	var jsonFile *os.File
	var wg sync.WaitGroup
	for _, zipFile := range zipReader.File {
		if jsonFile == nil && filepath.Ext(zipFile.Name) == "jsonl" {
			jsonFile, err = os.Open(zipFile.Name)
			if err != nil {
				panic(err)
			}
			continue
		}

		wg.Add(1)
		err := syscall.Mknod(zipFile.Name, syscall.S_IFIFO|0666, 0)
		if err != nil {
			panic(err)
		}

		go func(wg *sync.WaitGroup, zipFile *zip.File) {
			defer wg.Done()
			namedPipe, err := os.OpenFile(zipFile.Name, os.O_WRONLY, 0666)
			if err != nil {
				panic(err)
			}
			defer namedPipe.Close()

			fileAttachment, err := zipFile.Open()
			if err != nil {
				panic(err)
			}

			defer fileAttachment.Close()
			defer os.Remove(zipFile.Name)

			_, err = io.Copy(namedPipe, fileAttachment)
			if err != nil {
				panic(err)
			}

		}(&wg, zipFile)
	}

	// do the actual import.
	appErr, lineNumber := w.app.BulkImport(w.appContext, jsonFile, false, runtime.NumCPU())
	if appErr != nil {
		job.Data["line_number"] = strconv.Itoa(lineNumber)
		w.setJobError(job, appErr)
		return
	}

	wg.Wait()
	mlog.Info("Worker: Job is complete", mlog.String("worker", w.name), mlog.String("job_id", job.Id))
	w.setJobSuccess(job)
}

func (w *ImportProcessWorker) setJobSuccess(job *model.Job) {
	if err := w.app.Srv().Jobs.SetJobSuccess(job); err != nil {
		mlog.Error("Worker: Failed to set success for job", mlog.String("worker", w.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
		w.setJobError(job, err)
	}
}

func (w *ImportProcessWorker) setJobError(job *model.Job, appError *model.AppError) {
	if err := w.app.Srv().Jobs.SetJobError(job, appError); err != nil {
		mlog.Error("Worker: Failed to set job error", mlog.String("worker", w.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
	}
}
