package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	AppPort  string
	RedisURL string
}

var (
	config     *Config
	Configonce sync.Once
)

// GetConfig returns a singleton Config instance
func GetConfig() *Config {
	Configonce.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		viper.SetDefault("app_name", "Stock Scraper")
		viper.SetDefault("app_port", "3000")
		viper.SetDefault("redis_url", "redis://redis:6379")
		viper.AutomaticEnv()
		viper.BindEnv("redis_url", "REDIS_URL")

		// ✅ Env override support
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("Error reading config file: %v", err)
			}
		}

		config = &Config{
			AppName:  viper.GetString("app_name"),
			AppPort:  viper.GetString("app_port"),
			RedisURL: viper.GetString("redis_url"),
		}
	})

	return config
}
