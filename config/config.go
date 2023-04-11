package config

type DataStoreConfig struct {
	Target         string `mapstructure:"target"`
	URI            string `mapstructure:"uri"`
	Timeout        int    `mapstructure:"timeout"`
	DatabaseName   string `mapstructure:"databaseName"`
	CollectionName string `mapstructure:"collectionName"`
}

type BizObjConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	DataStores struct {
		// A map of data store configurations keyed by the data store name
		// e.g., "mongo" or "redis"
		DataStoreMap map[string]DataStoreConfig `mapstructure:"dataStore"`
	} `mapstructure:"dataStores"`
}

func New() *BizObjConfig {
	return &BizObjConfig{}
}
