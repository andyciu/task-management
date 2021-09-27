package models

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password *string
	Nickname *string
}
