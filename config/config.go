package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DataStoreConfig struct {
	Target string `mapstructure:"target"`
	URI string `mapstructure:"uri"`
	Timeout int `mapstructure:"timeout"`
	DatabaseName string `mapstructure:"databaseName"`
	CollectionName string `mapstructure:"collectionName"`
}

type BizObjConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	DataStores struct{
		// A map of data store configurations keyed by the data store name
		// e.g., "mongo" or "redis"
		DataStoreMap map[string]DataStoreConfig `mapstructure:"dataStore"`
		
	} `mapstructure:"dataStores"`
}


func LoadConfig(env string) (*BizObjConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf("config.%s.yaml", env)) // e.g., config.development.yaml or config.production.yaml
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
