// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/spf13/cobra"
)

var JobserverCmd = &cobra.Command{
	Use:   "jobserver",
	Short: "Start the Mattermost job server",
	Run:   jobserverCmdF,
}

func init() {
	JobserverCmd.Flags().Bool("nojobs", false, "Do not run jobs on this jobserver.")
	JobserverCmd.Flags().Bool("noschedule", false, "Do not schedule jobs from this jobserver.")

	RootCmd.AddCommand(JobserverCmd)
}

func jobserverCmdF(command *cobra.Command, args []string) {
	// Options
	noJobs, _ := command.Flags().GetBool("nojobs")
	noSchedule, _ := command.Flags().GetBool("noschedule")
	config, err := command.Flags().GetString("config")
	i18nOverride, _ := command.Flags().GetString("i18n-override")
	mailTemplatesOverride, _ := command.Flags().GetString("mail-templates-override")
	clientOverride, _ := command.Flags().GetString("i18n-override")

	// Initialize
	a, err := InitDBCommandContext(
		config,
		i18nOverride,
		mailTemplatesOverride,
		clientOverride,
	)

	if err != nil {
		panic(err.Error())
	}
	defer a.Shutdown()

	a.LoadLicense()

	// Run jobs
	mlog.Info("Starting Mattermost job server")
	if !noJobs {
		a.Jobs.StartWorkers()
	}
	if !noSchedule {
		a.Jobs.StartSchedulers()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Cleanup anything that isn't handled by a defer statement
	mlog.Info("Stopping Mattermost job server")

	a.Jobs.StopSchedulers()
	a.Jobs.StopWorkers()

	mlog.Info("Stopped Mattermost job server")
}
