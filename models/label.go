package models

type Label struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Task []*Task `gorm:"many2many:task_label_mapping;"`
}
