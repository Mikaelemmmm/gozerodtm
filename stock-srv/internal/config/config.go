package config

import (
	"github.com/tal-tech/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	//Consul consul.Conf //if consul use
	DB struct {
		DataSource string
	}
}
