package entities

import (
	"time"

	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	UserID uint
	Name   string
	User   User
	Task   []*Task `gorm:"many2many:task_label_mapping;"`
}

type Task struct {
	gorm.Model
	UserID      uint
	Title       string
	Description *string
	StartTime   *time.Time
	EndTime     *time.Time
	Priority    *int
	State       *int
	User        User
	Label       []*Label `gorm:"many2many:task_label_mapping;"`
}

type User struct {
	gorm.Model
	Username string
	Password *string
	Nickname *string
	AuthType uint //1-Local,2-GoogleOAuth
	Tasks    []*Task
	Labels   []*Label
}

type Userinfo_Google struct {
	gorm.Model
	UserID        uint
	UID           string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     *string
	FamilyName    *string
	Picture       *string
	Locale        *string
	User          User
}
