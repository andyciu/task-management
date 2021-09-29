package entities

import "time"

type Label struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Task []*Task `gorm:"many2many:task_label_mapping;"`
}

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

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password *string
	Nickname *string
}
