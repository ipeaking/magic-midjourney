package initialization

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DISCORD_USER_TOKEN string
	DISCORD_BOT_TOKEN  string
	DISCORD_SERVER_ID  string
	DISCORD_CHANNEL_ID string
	CB_URL             string
	QiNiuConfig        *qiNiuConfig
}

type qiNiuConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
}

var config *Config

func LoadConfig(cfg string) (*Config, error) {
	viper.SetConfigFile(cfg)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	config = &Config{
		DISCORD_USER_TOKEN: getViperStringValue("DISCORD_USER_TOKEN"),
		DISCORD_BOT_TOKEN:  getViperStringValue("DISCORD_BOT_TOKEN"),
		DISCORD_SERVER_ID:  getViperStringValue("DISCORD_SERVER_ID"),
		DISCORD_CHANNEL_ID: getViperStringValue("DISCORD_CHANNEL_ID"),
		CB_URL:             getViperStringValue("CB_URL"),
		QiNiuConfig: &qiNiuConfig{
			AccessKey: getViperStringValue("QINIU_ACCESS_KEY"),
			SecretKey: getViperStringValue("QINIU_SECRET_KEY"),
			Bucket:    getViperStringValue("QINIU_BUCKET"),
		},
	}

	return config, nil
}

func GetConfig() *Config {
	return config
}

func getViperStringValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s MUST be provided in environment or config.yaml file", key))
	}
	return value
}
