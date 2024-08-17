package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func exportConfig(path string) error {
	var configPath string
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if path == "" {
		configPath = fmt.Sprintf("config.%s.yaml", os.Getenv("APP_ENV"))
	} else {
		configPath = path
	}

	viper.SetConfigName(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig Parse config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(""); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("Unable to decode config file into struct, %v \n", err)
		return nil, err
	}

	return &c, nil
}

// ParseMockConfig Parse config file
func ParseMockConfig(path string) (*Config, error) {
	if err := exportConfig(path); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("Unable to decode config file into struct, %v \n", err)
		return nil, err
	}

	return &c, nil
}
