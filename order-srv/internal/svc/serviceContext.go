package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"gozerodtm/order-srv/internal/config"
	"gozerodtm/order-srv/internal/model"
)

type ServiceContext struct {
	Config config.Config
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		OrderModel:  model.NewOrderModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
