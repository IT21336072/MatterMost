// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (api *API) InitReports() {
	api.BaseRoutes.Reports.Handle("/users", api.APISessionRequired(getUsersForReporting)).Methods("GET")
	api.BaseRoutes.Reports.Handle("/users/export", api.APISessionRequired(startUsersBatchExport)).Methods("POST")

	api.BaseRoutes.Reports.Handle("/export/{report_id:[A-Za-z0-9]+}", api.APISessionRequired(retrieveBatchReportFile)).Methods("GET")
}

func getUsersForReporting(c *Context, w http.ResponseWriter, r *http.Request) {
	if !(c.IsSystemAdmin() && c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementUsers)) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementUsers)
		return
	}

	sortColumn := "Username"
	if r.URL.Query().Get("sort_column") != "" {
		sortColumn = r.URL.Query().Get("sort_column")
	}

	direction := "down"
	if r.URL.Query().Get("direction") == "up" {
		direction = "up"
	}

	pageSize := 50
	if pageSizeStr, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64); err == nil {
		pageSize = int(pageSizeStr)
	}

	teamFilter := r.URL.Query().Get("team_filter")
	if !(teamFilter == "" || model.IsValidId(teamFilter)) {
		c.Err = model.NewAppError("getUsersForReporting", "api.getUsersForReporting.invalid_team_filter", nil, "", http.StatusBadRequest)
		return
	}

	hideActive := r.URL.Query().Get("hide_active") == "true"
	hideInactive := r.URL.Query().Get("hide_inactive") == "true"
	if hideActive && hideInactive {
		c.Err = model.NewAppError("getUsersForReporting", "api.getUsersForReporting.invalid_active_filter", nil, "", http.StatusBadRequest)
		return
	}

	options := &model.UserReportOptions{
		ReportingBaseOptions: model.ReportingBaseOptions{
			Direction:       direction,
			SortColumn:      sortColumn,
			SortDesc:        r.URL.Query().Get("sort_direction") == "desc",
			PageSize:        pageSize,
			FromColumnValue: r.URL.Query().Get("from_column_value"),
			FromId:          r.URL.Query().Get("from_id"),
			DateRange:       r.URL.Query().Get("date_range"),
		},
		Team:         teamFilter,
		Role:         r.URL.Query().Get("role_filter"),
		HasNoTeam:    r.URL.Query().Get("has_no_team") == "true",
		HideActive:   hideActive,
		HideInactive: hideInactive,
	}
	options.PopulateDateRange(time.Now())

	// Don't allow fetching more than 100 users at a time from the normal query endpoint
	if options.PageSize <= 0 || options.PageSize > model.ReportingMaxPageSize {
		c.Err = model.NewAppError("getUsersForReporting", "api.getUsersForReporting.invalid_page_size", nil, "", http.StatusBadRequest)
		return
	}

	userReports, err := c.App.GetUsersForReporting(options)
	if err != nil {
		c.Err = err
		return
	}

	if jsonErr := json.NewEncoder(w).Encode(userReports); jsonErr != nil {
		c.Logger.Warn("Error writing response", mlog.Err(jsonErr))
	}
}

func startUsersBatchExport(c *Context, w http.ResponseWriter, r *http.Request) {
	if !(c.IsSystemAdmin() && c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementUsers)) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementUsers)
		return
	}

	if err := c.App.StartUsersBatchExport(c.AppContext, r.URL.Query().Get("date_range")); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func retrieveBatchReportFile(c *Context, w http.ResponseWriter, r *http.Request) {
	if !(c.IsSystemAdmin() && c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementUsers)) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementUsers)
		return
	}

	reportId := c.Params.ReportId
	if reportId == "" || !model.IsValidId(reportId) {
		c.Err = model.NewAppError("retrieveBatchReportFile", "api.retrieveBatchReportFile.invalid_report_id", nil, "", http.StatusBadRequest)
	}

	format := r.URL.Query().Get("format")
	// TODO: Validate with more types
	if format != "csv" {
		c.Err = model.NewAppError("retrieveBatchReportFile", "api.retrieveBatchReportFile.invalid_format", nil, "", http.StatusBadRequest)
	}

	file, name, err := c.App.RetrieveBatchReport(reportId, format)
	if err != nil {
		c.Err = err
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/csv") // TODO: other formats
	http.ServeContent(w, r, name, time.Time{}, file)
}
