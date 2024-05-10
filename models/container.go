package models

import "time"

type Container struct {
	ID           uint   `gorm:"PrimaryKey"`
	Name         string `gorm:"index:idx_name,unique"`
	Namespace    string `gorm:"index:idx_name,unique"`
	CpuRequest   string
	CpuLimit     string
	MemReq       string
	MemLimit     string
	Replicas     int32
	ImageVersion string
	CreateAt     time.Time
}

func (Container) TableName() string {
	return "container"
}
