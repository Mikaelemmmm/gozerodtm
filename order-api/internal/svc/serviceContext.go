package svc

import (
	"github.com/tal-tech/go-zero/zrpc"
	"gozerodtm/order-api/internal/config"
	"gozerodtm/order-srv/order"
	"gozerodtm/stock-srv/stock"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc order.Order
	StockRpc stock.Stock
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc: order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		StockRpc: stock.NewStock(zrpc.MustNewClient(c.StockRpcConf)),
	}
}
