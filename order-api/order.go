package main

import (
	"flag"
	"fmt"

	// 导入gozero驱动
	_ "github.com/yedf/driver-gozero"

	"gozerodtm/order-api/internal/config"
	"gozerodtm/order-api/internal/handler"
	"gozerodtm/order-api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
	//_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"  //if consul use
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
