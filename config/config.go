package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func New(path string) (config Config, err error) {
	var configName string
	switch os.Getenv("APP_ENV") {
	case "test":
		configName = "test"
	case "dev":
		configName = "dev"
	case "prod":
		configName = "prod"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
