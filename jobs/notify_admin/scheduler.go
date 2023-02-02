// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package notify_admin

import (
	"strconv"
	"time"

	"github.com/mattermost/mattermost-server/v6/jobs"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

const schedFreq = 1 * time.Minute

func MakeScheduler(jobServer *jobs.JobServer, license *model.License, jobType string) model.Scheduler {
	isEnabled := func(cfg *model.Config) bool {
		installPluginRequestJobType := jobType == model.JobTypeInstallPluginNotifyAdmin
		enabled := (license != nil && *license.Features.Cloud) || installPluginRequestJobType
		mlog.Debug("Scheduler: isEnabled: "+strconv.FormatBool(enabled), mlog.String("scheduler", jobType))
		return enabled
	}
	return jobs.NewPeriodicScheduler(jobServer, jobType, schedFreq, isEnabled)

}
