package config

import "github.com/spf13/viper"

type Config struct {
	TelegramToken string `mapstructure:""`
	DatabaseUrl   string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
