package models

import (
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/interfaces"
)

// Constructor Options for IModel types that allows you to set the logging level
// at construction
func WithLoggingLevel(level string) func(interfaces.IModel) {
	return func(obj interfaces.IModel) {
		obj.SetLoggingLevel(level)
	}
}

// Constructor Options for IModel types that allows you to set the BizObjConfig
// at construction
func WithBizObjConfig(config config.BizObjConfig) func(interfaces.IModel) {
	return func(obj interfaces.IModel) {
		obj.SetBizObjConfig(config)
	}
}

// Constructor Options for IModel types that allows you to set the DataStoreConfig
// at construction
func WithDataStoresConfig(config config.DataStoresConfig) func(interfaces.IModel) {
	return func(obj interfaces.IModel) {
		obj.SetDataStoresConfig(config)
	}
}
