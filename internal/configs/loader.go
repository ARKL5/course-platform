package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 結構體定義了應用程式的所有配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// MySQLConfig MySQL 資料庫配置
type MySQLConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoadConfig 讀取並解析配置檔案
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// 讀取配置檔案
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("讀取配置檔案失敗: %v", err)
	}

	var config Config
	// 將配置解析到結構體中
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置檔案失敗: %v", err)
	}

	return &config, nil
}
