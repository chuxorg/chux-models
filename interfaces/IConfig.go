package interfaces

import "github.com/csailer/chux-bizobj/config"

type IConfig interface {
	LoadConfig() (*config.BizObjConfig, error)
}
