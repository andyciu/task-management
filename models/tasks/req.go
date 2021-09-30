package tasks

import (
	"encoding/json"
	"time"
)

type TaskListReq struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Labels      []json.Number `json:"labels,omitempty"`
}

type TaskCreateReq struct {
	Title       string        `json:"title" binding:"required"`
	Description *string       `json:"description,omitempty"`
	StartTime   *time.Time    `json:"start_time,omitempty"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	Priority    *json.Number  `json:"priority,omitempty"`
	State       *json.Number  `json:"state,omitempty"`
	Labels      []json.Number `json:"labels,omitempty"`
}

type TaskUpdateReq struct {
	ID          json.Number   `json:"id" binding:"required"`
	Title       string        `json:"title" binding:"required"`
	Description *string       `json:"description,omitempty"`
	StartTime   *time.Time    `json:"start_time,omitempty"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	Priority    *json.Number  `json:"priority,omitempty"`
	State       *json.Number  `json:"state,omitempty"`
	Labels      []json.Number `json:"labels,omitempty"`
}

type TaskDeleteReq struct {
	ID json.Number `json:"id" binding:"required"`
}
