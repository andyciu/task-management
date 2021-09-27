package models

import "time"

type Task struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Title       string
	Description *string
	StartTime   time.Time
	EndTime     time.Time
	Priority    *int
	State       *int
	User        User
	Label       []*Label `gorm:"many2many:task_label_mapping;"`
}
