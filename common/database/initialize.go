package database

import (
	"strconv"
	"workassistant/common"
	"workassistant/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func RunMigrations() error {
	var err error
	setting := common.GetConfig()
	dsn := setting.GetString("settings.database.source")
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
