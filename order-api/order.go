package main

import (
	"flag"
	"fmt"
	"github.com/yedf/dtmcli/dtmimp"
	"github.com/yedf/dtmdriver"

	"gozerodtm/order-api/internal/config"
	"gozerodtm/order-api/internal/handler"
	"gozerodtm/order-api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
	//导入驱动
	_ "github.com/yedf/dtmdriver-gozero"
)


var configFile = flag.String("f", "/Users/mikael/Developer/goenv/gozerodtm/order-api/etc/order.yaml", "the config file")

func main() {

	flag.Parse()

	// 使用dtm的客户端dtmgrpc之前，需要执行下面这行调用，告知dtmgrpc使用gozero的驱动来如何处理gozero的url
	err := dtmdriver.Use("dtm-driver-gozero")
	dtmimp.FatalIfError(err)

	var c config.Config
	conf.MustLoad(*configFile, &c)


	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
