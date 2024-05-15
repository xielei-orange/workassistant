// config.go
package common

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	Config     *viper.Viper
	configOnce sync.Once
)

func initConfig() {
	configOnce.Do(func() {
		Config = viper.New()
		Config.SetConfigName("settings-dev")
		Config.SetConfigType("yaml")
		Config.AddConfigPath("/data/workassistant/config")
		Config.AddConfigPath("config/")

		// 读取配置文件
		if err := Config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// 配置文件未找到，但可能不是错误，可以根据需要处理
				fmt.Println("Warning: Using default config because not found")
			} else {
				// 发生其他错误，可能是配置文件格式错误等
				fmt.Printf("Fatal error config file: %s \n", err)
			}
		}
	})
}

// InitConfig 初始化配置
func InitConfig() {
	initConfig()
}

// GetConfig 返回全局配置对象
func GetConfig() *viper.Viper {
	initConfig()
	return Config
}
