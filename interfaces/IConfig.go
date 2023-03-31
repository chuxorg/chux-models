package interfaces

import "github.com/chuxorg/chux-models/config"

type IConfig interface {
	LoadConfig() (*config.BizObjConfig, error)
}
