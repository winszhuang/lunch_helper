package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ApiUrl                    string `mapstructure:"API_URL"`
	DBDriver                  string `mapstructure:"DB_DRIVER"`
	DBSource                  string `mapstructure:"DB_SOURCE"`
	LineBotChannelSecret      string `mapstructure:"LINEBOT_CHANNEL_SECRET"`
	LineBotChannelAccessToken string `mapstructure:"LINEBOT_CHANNEL_ACCESS_TOKEN"`
	LineBotEndpoint           string `mapstructure:"LINEBOT_ENDPOINT"`
}

func New(path, env string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
