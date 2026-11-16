package taskService

import "gorm.io/gorm"

type Task struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
