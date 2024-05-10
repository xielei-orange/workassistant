package models

import (
	"gorm.io/gorm"
)

type ChangeLog struct {
	gorm.Model
	Namespace         string
	Verb              string //操作类型
	Status            int
	ObjectRefResource string
	ObjectRefName     string
	SourceIPs         string
	User              string
}

func (ChangeLog) TableName() string {
	return "changelog"
}
