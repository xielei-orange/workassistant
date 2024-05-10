package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"tools/models"
)

var (
	DB *gorm.DB
)

func RunMigrations() error {
	var err error
	dsn := "root:x2eracom@tcp(10.40.64.32:3306)/devops?charset=utf8mb4&parseTime=True&loc=Local"
	config := gorm.Config{}
	DB, err = gorm.Open(mysql.Open(dsn), &config)
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&models.Container{}, &models.ChangeLog{})
	if err != nil {
		return err
	}
	return nil
}

func insert(container models.Container) (string, error) {
	resp := DB.Create(&container)
	if resp.Error != nil {
		return "数据插入失败", resp.Error
	}
	return strconv.FormatInt(resp.RowsAffected, 10), nil
}
