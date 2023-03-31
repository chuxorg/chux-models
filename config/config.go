package config

import (
	"fmt"

	"github.com/csailer/chux-mongo/db"
	"github.com/spf13/viper"
)

type BizObjConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	MongoDB db.MongoConfig `mapstructure:"mongodb"`
}

func LoadConfig(env string) (*BizObjConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf("config.%s", env)) // e.g., config.development.yaml or config.production.yaml
	viper.AddConfigPath(".")                           // Look for config files in the current directory
	viper.AddConfigPath("./config")                    // Look for config files in the config directory
	viper.AddConfigPath("../config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	var cfg BizObjConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	return &cfg, nil
}
