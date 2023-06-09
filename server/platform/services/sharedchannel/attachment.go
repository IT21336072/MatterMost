// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app/request"
	"github.com/mattermost/mattermost/server/v8/platform/services/remotecluster"
)

// postsToAttachments returns the file attachments for a slice of posts that need to be synchronized.
func (scs *Service) shouldSyncAttachment(fi *model.FileInfo, rc *model.RemoteCluster) bool {
	sca, err := scs.server.GetStore().SharedChannel().GetAttachment(fi.Id, rc.RemoteId)
	if err != nil {
		if _, ok := err.(errNotFound); !ok {
			scs.server.Log().Log(mlog.LvlSharedChannelServiceError, "error fetching shared channel attachment",
				mlog.String("file_id", fi.Id),
				mlog.String("remote_id", rc.RemoteId),
				mlog.Err(err),
			)
		}
		// no record so sync is needed
		return true
	}

	return sca.LastSyncAt < fi.UpdateAt
}

// sendAttachmentForRemote asynchronously sends a file attachment to a remote cluster.
func (scs *Service) sendAttachmentForRemote(fi *model.FileInfo, post *model.Post, rc *model.RemoteCluster) error {
	rcs := scs.server.GetRemoteClusterService()
	if rcs == nil {
		return fmt.Errorf("cannot update remote cluster for remote id %s; Remote Cluster Service not enabled", rc.RemoteId)
	}

	us := &model.UploadSession{
		Id:        model.NewId(),
		Type:      model.UploadTypeAttachment,
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Filename:  fi.Name,
		FileSize:  fi.Size,
		RemoteId:  rc.RemoteId,
		ReqFileId: fi.Id,
	}

	payload, err := json.Marshal(us)
	if err != nil {
		return err
	}

	msg := model.NewRemoteClusterMsg(TopicUploadCreate, payload)

	ctx, cancel := context.WithTimeout(context.Background(), remotecluster.SendTimeout)
	defer cancel()

	var usResp model.UploadSession
	var respErr error
	var wg sync.WaitGroup
	wg.Add(1)

	// creating the upload session on the remote server needs to be done synchronously.
	err = rcs.SendMsg(ctx, msg, rc, func(msg model.RemoteClusterMsg, rc *model.RemoteCluster, resp *remotecluster.Response, err error) {
		defer wg.Done()
		if err != nil {
			respErr = err
			return
		}
		if !resp.IsSuccess() {
			respErr = errors.New(resp.Err)
			return
		}
		respErr = json.Unmarshal(resp.Payload, &usResp)
	})

	if err != nil {
		return fmt.Errorf("error sending create upload session to remote %s for post %s: %w", rc.RemoteId, post.Id, err)
	}

	wg.Wait()

	if respErr != nil {
		return fmt.Errorf("invalid create upload session response for remote %s and post %s: %w", rc.RemoteId, post.Id, respErr)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), remotecluster.SendFileTimeout)
	defer cancel2()

	return rcs.SendFile(ctx2, &usResp, fi, rc, scs.app, func(us *model.UploadSession, rc *model.RemoteCluster, resp *remotecluster.Response, err error) {
		if err != nil {
			return // this means the response could not be parsed; already logged
		}

		if !resp.IsSuccess() {
			scs.server.Log().Log(mlog.LvlSharedChannelServiceError, "send file failed",
				mlog.String("remote", rc.DisplayName),
				mlog.String("uploadId", usResp.Id),
				mlog.String("err", resp.Err),
			)
			return
		}

		// response payload should be a model.FileInfo.
		var fi model.FileInfo
		if err2 := json.Unmarshal(resp.Payload, &fi); err2 != nil {
			scs.server.Log().Log(mlog.LvlSharedChannelServiceError, "invalid file info response after send file",
				mlog.String("remote", rc.DisplayName),
				mlog.String("uploadId", usResp.Id),
				mlog.Err(err2),
			)
			return
		}

		// save file attachment record in SharedChannelAttachments table
		sca := &model.SharedChannelAttachment{
			FileId:   fi.Id,
			RemoteId: rc.RemoteId,
		}
		if _, err2 := scs.server.GetStore().SharedChannel().UpsertAttachment(sca); err2 != nil {
			scs.server.Log().Log(mlog.LvlSharedChannelServiceError, "error saving SharedChannelAttachment",
				mlog.String("remote", rc.DisplayName),
				mlog.String("uploadId", usResp.Id),
				mlog.Err(err2),
			)
			return
		}

		scs.server.Log().Log(mlog.LvlSharedChannelServiceDebug, "send file successful",
			mlog.String("remote", rc.DisplayName),
			mlog.String("uploadId", usResp.Id),
		)
	})
}

// onReceiveUploadCreate is called when a message requesting to create an upload session is received.  An upload session is
// created and the id returned in the response.
func (scs *Service) onReceiveUploadCreate(msg model.RemoteClusterMsg, rc *model.RemoteCluster, response *remotecluster.Response) error {
	var us model.UploadSession

	if err := json.Unmarshal(msg.Payload, &us); err != nil {
		return fmt.Errorf("invalid upload session request: %w", err)
	}

	// make sure channel is shared for the remote sender
	if _, err := scs.server.GetStore().SharedChannel().GetRemoteByIds(us.ChannelId, rc.RemoteId); err != nil {
		return fmt.Errorf("could not validate upload session for remote: %w", err)
	}

	us.RemoteId = rc.RemoteId // don't let remotes try to impersonate each other

	// create upload session.
	usSaved, appErr := scs.app.CreateUploadSession(request.EmptyContext(scs.server.Log()), &us)
	if appErr != nil {
		return appErr
	}

	response.SetPayload(usSaved)
	return nil
}
