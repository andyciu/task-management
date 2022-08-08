package tasks

import (
	"time"
)

type TaskListRes struct {
	ID          int        `json:"id"`
	Num         int        `json:"num"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Priority    *int       `json:"priority"`
	State       *int       `json:"state"`
	Labels      []int      `json:"labels"`
}
