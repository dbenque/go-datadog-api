/*
 * Datadog API for Go
 *
 * Please see the included LICENSE file for licensing information.
 *
 * Copyright 2013 by authors and contributors.
 */

package datadog

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Period struct {
	Seconds *json.Number `json:"seconds,omitempty"`
	Text    *string      `json:"text,omitempty"`
	Value   *string      `json:"value,omitempty"`
	Name    *string      `json:"name,omitempty"`
	Unit    *string      `json:"unit,omitempty"`
}

type LogSet struct {
	ID   *json.Number `json:"id,omitempty"`
	Name *string      `json:"name,omitempty"`
}

type TimeRange struct {
	To   *json.Number `json:"to,omitempty"`
	From *json.Number `json:"from,omitempty"`
	Live *bool        `json:"live,omitempty"`
}

type QueryConfig struct {
	LogSet        *LogSet    `json:"logset,omitempty"`
	TimeRange     *TimeRange `json:"timeRange,omitempty"`
	QueryString   *string    `json:"queryString,omitempty"`
	QueryIsFailed *bool      `json:"queryIsFailed,omitempty"`
}

type ThresholdCount struct {
	Ok               *json.Number `json:"ok,omitempty"`
	Critical         *json.Number `json:"critical,omitempty"`
	Warning          *json.Number `json:"warning,omitempty"`
	Unknown          *json.Number `json:"unknown,omitempty"`
	CriticalRecovery *json.Number `json:"critical_recovery,omitempty"`
	WarningRecovery  *json.Number `json:"warning_recovery,omitempty"`
	Period           *Period      `json:"period,omitempty"`
	TimeAggregator   *string      `json:"timeAggregator,omitempty"`
}

type ThresholdWindows struct {
	RecoveryWindow *string `json:"recovery_window,omitempty"`
	TriggerWindow  *string `json:"trigger_window,omitempty"`
}

type NoDataTimeframe int

func (tf *NoDataTimeframe) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "false" || s == "null" {
		*tf = 0
	} else {
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		*tf = NoDataTimeframe(i)
	}
	return nil
}

type Options struct {
	NoDataTimeframe   NoDataTimeframe   `json:"no_data_timeframe,omitempty"`
	NotifyAudit       *bool             `json:"notify_audit,omitempty"`
	NotifyNoData      *bool             `json:"notify_no_data,omitempty"`
	RenotifyInterval  *int              `json:"renotify_interval,omitempty"`
	NewHostDelay      *int              `json:"new_host_delay,omitempty"`
	EvaluationDelay   *int              `json:"evaluation_delay,omitempty"`
	Silenced          map[string]int    `json:"silenced,omitempty"`
	TimeoutH          *int              `json:"timeout_h,omitempty"`
	EscalationMessage *string           `json:"escalation_message,omitempty"`
	Thresholds        *ThresholdCount   `json:"thresholds,omitempty"`
	ThresholdWindows  *ThresholdWindows `json:"threshold_windows,omitempty"`
	IncludeTags       *bool             `json:"include_tags,omitempty"`
	RequireFullWindow *bool             `json:"require_full_window,omitempty"`
	Locked            *bool             `json:"locked,omitempty"`
	EnableLogsSample  *bool             `json:"enable_logs_sample,omitempty"`
	QueryConfig       *QueryConfig      `json:"queryConfig,omitempty"`
}

type TriggeringValue struct {
	FromTs *int `json:"from_ts,omitempty"`
	ToTs   *int `json:"to_ts,omitempty"`
	Value  *int `json:"value,omitempty"`
}

type GroupData struct {
	LastNoDataTs    *int             `json:"last_nodata_ts,omitempty"`
	LastNotifiedTs  *int             `json:"last_notified_ts,omitempty"`
	LastResolvedTs  *int             `json:"last_resolved_ts,omitempty"`
	LastTriggeredTs *int             `json:"last_triggered_ts,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Status          *string          `json:"status,omitempty"`
	TriggeringValue *TriggeringValue `json:"triggering_value,omitempty"`
}

type State struct {
	Groups map[string]GroupData `json:"groups,omitempty"`
}

// Monitor allows watching a metric or check that you care about,
// notifying your team when some defined threshold is exceeded
type Monitor struct {
	Creator              *Creator `json:"creator,omitempty"`
	Id                   *int     `json:"id,omitempty"`
	Type                 *string  `json:"type,omitempty"`
	Query                *string  `json:"query,omitempty"`
	Name                 *string  `json:"name,omitempty"`
	Message              *string  `json:"message,omitempty"`
	OverallState         *string  `json:"overall_state,omitempty"`
	OverallStateModified *string  `json:"overall_state_modified,omitempty"`
	Tags                 []string `json:"tags"`
	Options              *Options `json:"options,omitempty"`
	State                State    `json:"state,omitempty"`
}

// Creator contains the creator of the monitor
type Creator struct {
	Email  *string `json:"email,omitempty"`
	Handle *string `json:"handle,omitempty"`
	Id     *int    `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
}

// MuteMonitorScope specifies which scope to mute and when to end the mute
type MuteMonitorScope struct {
	Scope *string `json:"scope,omitempty"`
	End   *int    `json:"end,omitempty"`
}

// UnmuteMonitorScopes specifies which scope(s) to unmute
type UnmuteMonitorScopes struct {
	Scope     *string `json:"scope,omitempty"`
	AllScopes *bool   `json:"all_scopes,omitempty"`
}

// reqMonitors receives a slice of all monitors
type reqMonitors struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

// CreateMonitor adds a new monitor to the system. This returns a pointer to a
// monitor so you can pass that to UpdateMonitor later if needed
func (client *Client) CreateMonitor(monitor *Monitor) (*Monitor, error) {
	var out Monitor
	// TODO: is this more pretty of frowned upon?
	if err := client.doJsonRequest("POST", "/v1/monitor", monitor, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMonitor takes a monitor that was previously retrieved through some method
// and sends it back to the server
func (client *Client) UpdateMonitor(monitor *Monitor) error {
	return client.doJsonRequest("PUT", fmt.Sprintf("/v1/monitor/%d", *monitor.Id),
		monitor, nil)
}

// GetMonitor retrieves a monitor by identifier
func (client *Client) GetMonitor(id int) (*Monitor, error) {
	var out Monitor
	if err := client.doJsonRequest("GET", fmt.Sprintf("/v1/monitor/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMonitorsByName retrieves monitors by name
func (self *Client) GetMonitorsByName(name string) ([]Monitor, error) {
	var out reqMonitors
	query, err := url.ParseQuery(fmt.Sprintf("name=%v", name))
	if err != nil {
		return nil, err
	}

	err = self.doJsonRequest("GET", fmt.Sprintf("/v1/monitor?%v", query.Encode()), nil, &out.Monitors)
	if err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// GetMonitorsByTags retrieves monitors by a slice of tags
func (self *Client) GetMonitorsByTags(monitorTags []string) ([]Monitor, error) {
	var out reqMonitors
	query, err := url.ParseQuery(fmt.Sprintf("monitor_tags=%v", strings.Join(monitorTags, ",")))
	if err != nil {
		return nil, err
	}

	err = self.doJsonRequest("GET", fmt.Sprintf("/v1/monitor?%v", query.Encode()), nil, &out.Monitors)
	if err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// GetMonitorsWithGroupsByTags retrieves monitors by tags, monitor_tags, groupState
func (self *Client) GetMonitorsWithGroupsByTags(monitorTags []string, tags []string, groupState []string, withDowntimes bool) ([]Monitor, error) {
	var out reqMonitors

	queryStr := ""
	if len(monitorTags) > 0 {
		queryStr = fmt.Sprintf("monitor_tags=%v", strings.Join(monitorTags, ","))
	}
	if len(tags) > 0 {
		if queryStr != "" {
			queryStr += "&"
		}
		queryStr += fmt.Sprintf("tags=%v", strings.Join(tags, ","))
	}
	if len(groupState) > 0 {
		if queryStr != "" {
			queryStr += "&"
		}
		queryStr += fmt.Sprintf("group_states=%v", strings.Join(groupState, ","))
	}

	if queryStr != "" {
		queryStr += "&"
	}
	if withDowntimes {
		queryStr += "with_downtimes=true"
	} else {
		queryStr += "with_downtimes=false"
	}

	query, err := url.ParseQuery(queryStr)
	if err != nil {
		return nil, err
	}

	err = self.doJsonRequest("GET", fmt.Sprintf("/v1/monitor?%v", query.Encode()), nil, &out.Monitors)
	if err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// DeleteMonitor removes a monitor from the system
func (client *Client) DeleteMonitor(id int) error {
	return client.doJsonRequest("DELETE", fmt.Sprintf("/v1/monitor/%d", id),
		nil, nil)
}

// GetMonitors returns a slice of all monitors
func (client *Client) GetMonitors() ([]Monitor, error) {
	var out reqMonitors
	if err := client.doJsonRequest("GET", "/v1/monitor", nil, &out.Monitors); err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// MuteMonitors turns off monitoring notifications
func (client *Client) MuteMonitors() error {
	return client.doJsonRequest("POST", "/v1/monitor/mute_all", nil, nil)
}

// UnmuteMonitors turns on monitoring notifications
func (client *Client) UnmuteMonitors() error {
	return client.doJsonRequest("POST", "/v1/monitor/unmute_all", nil, nil)
}

// MuteMonitor turns off monitoring notifications for a monitor
func (client *Client) MuteMonitor(id int) error {
	return client.doJsonRequest("POST", fmt.Sprintf("/v1/monitor/%d/mute", id), nil, nil)
}

// MuteMonitorScope turns off monitoring notifications for a monitor for a given scope
func (client *Client) MuteMonitorScope(id int, muteMonitorScope *MuteMonitorScope) error {
	return client.doJsonRequest("POST", fmt.Sprintf("/v1/monitor/%d/mute", id), muteMonitorScope, nil)
}

// UnmuteMonitor turns on monitoring notifications for a monitor
func (client *Client) UnmuteMonitor(id int) error {
	return client.doJsonRequest("POST", fmt.Sprintf("/v1/monitor/%d/unmute", id), nil, nil)
}

// UnmuteMonitorScopes is similar to UnmuteMonitor, but provides finer-grained control to unmuting
func (client *Client) UnmuteMonitorScopes(id int, unmuteMonitorScopes *UnmuteMonitorScopes) error {
	return client.doJsonRequest("POST", fmt.Sprintf("/v1/monitor/%d/unmute", id), unmuteMonitorScopes, nil)
}
