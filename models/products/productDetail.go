package products

import (
	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
)

type ProductDetail struct {
	Brand string `bson:"brand,omitempty"`
	URL   string `bson:"url,omitempty"`
}

var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func New(options ...func(*ProductDetail)) *ProductDetail {

	_cfg = config.New()
	productDetail := &ProductDetail{}
	for _, option := range options {
		option(productDetail)
	}

	mongoDB = db.New(
		db.WithURI(product.GetURI()),
		db.WithDatabaseName(product.GetDatabaseName()),
		db.WithCollectionName(product.GetCollectionName()),
		db.WithTimeout(float64(_cfg.DataStores.DataStoreMap["mongo"].Timeout)),
	)

	product.isNew = true
	product.isDeleted = false
	product.isDirty = false
	return product
}
