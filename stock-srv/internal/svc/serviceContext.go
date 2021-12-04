package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"gozerodtm/stock-srv/internal/config"
	"gozerodtm/stock-srv/internal/model"
)

type ServiceContext struct {
	Config config.Config
	StockModel model.StockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		StockModel:  model.NewStockModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
