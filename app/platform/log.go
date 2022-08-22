// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package platform

import (
	"context"
	"errors"
	"time"

	"github.com/mattermost/mattermost-server/v6/config"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

func (ps *PlatformService) ReconfigureLogger() error {
	return ps.initLogging()
}

// initLogging initializes and configures the logger(s). This may be called more than once.
func (ps *PlatformService) initLogging() error {
	var err error
	// create the app logger if needed
	if ps.logger == nil {
		ps.logger, err = mlog.NewLogger()
		if err != nil {
			return err
		}
	}

	// create notification logger if needed
	if ps.notificationsLogger == nil {
		l, err := mlog.NewLogger()
		if err != nil {
			return err
		}
		ps.notificationsLogger = l.With(mlog.String("logSource", "notifications"))
	}

	if err := ps.ConfigureLogger("logging", ps.logger, &ps.Config().LogSettings, config.GetLogFileLocation); err != nil {
		// if the config is locked then a unit test has already configured and locked the logger; not an error.
		if !errors.Is(err, mlog.ErrConfigurationLock) {
			// revert to default logger if the config is invalid
			mlog.InitGlobalLogger(nil)
			return err
		}
	}

	// Redirect default Go logger to app logger.
	ps.logger.RedirectStdLog(mlog.LvlStdLog)

	// Use the app logger as the global logger (eventually remove all instances of global logging).
	mlog.InitGlobalLogger(ps.logger)

	notificationLogSettings := config.GetLogSettingsFromNotificationsLogSettings(&ps.Config().NotificationLogSettings)
	if err := ps.ConfigureLogger("notification logging", ps.notificationsLogger, notificationLogSettings, config.GetNotificationsLogFileLocation); err != nil {
		if !errors.Is(err, mlog.ErrConfigurationLock) {
			mlog.Error("Error configuring notification logger", mlog.Err(err))
			return err
		}
	}
	return nil
}

func (ps *PlatformService) Logger() *mlog.Logger {
	return ps.logger
}

func (ps *PlatformService) NotificationsLogger() *mlog.Logger {
	return ps.notificationsLogger
}

func (ps *PlatformService) EnableLoggingMetrics() {
	if ps.metrics == nil || ps.metrics.metricsImpl == nil {
		return
	}

	ps.logger.SetMetricsCollector(ps.metrics.metricsImpl.GetLoggerMetricsCollector(), mlog.DefaultMetricsUpdateFreqMillis)

	// logging config needs to be reloaded when metrics collector is added or changed.
	if err := ps.initLogging(); err != nil {
		mlog.Error("Error re-configuring logging for metrics")
		return
	}

	mlog.Debug("Logging metrics enabled")
}

// RemoveUnlicensedLogTargets removes any unlicensed log target types.
func (ps *PlatformService) RemoveUnlicensedLogTargets(license *model.License) {
	if license != nil && *license.Features.AdvancedLogging {
		// advanced logging enabled via license; no need to remove any targets
		return
	}

	timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelCtx()

	ps.logger.RemoveTargets(timeoutCtx, func(ti mlog.TargetInfo) bool {
		return ti.Type != "*targets.Writer" && ti.Type != "*targets.File"
	})

	ps.notificationsLogger.RemoveTargets(timeoutCtx, func(ti mlog.TargetInfo) bool {
		return ti.Type != "*targets.Writer" && ti.Type != "*targets.File"
	})
}
