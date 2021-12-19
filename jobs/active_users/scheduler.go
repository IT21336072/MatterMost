// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package active_users

import (
	"github.com/mattermost/mattermost-server/v6/jobs"
	"github.com/mattermost/mattermost-server/v6/model"
)

const (
	SchedFreqMinutes = 10
)

func (m *ActiveUsersJobInterfaceImpl) MakeScheduler() model.Scheduler {
	isEnabled := func(cfg *model.Config) bool {
		return *cfg.MetricsSettings.Enable
	}
	return jobs.NewPeridicScheduler(m.App.Srv().Jobs, model.JobTypeActiveUsers, SchedFreqMinutes, isEnabled)
}
