package main

import "github.com/spf13/viper"

const (
	mnemonicField = "mnemonic"
)

type Config struct {
	Mnemonic string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		Mnemonic: viper.GetString(mnemonicField),
	}, nil
}
