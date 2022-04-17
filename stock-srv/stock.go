package main

import (
	"flag"
	"fmt"
	"gozerodtm/stock-srv/internal/config"
	"gozerodtm/stock-srv/internal/server"
	"gozerodtm/stock-srv/internal/svc"
	"gozerodtm/stock-srv/pb"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/zrpc"
	//_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul" //if consul use
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/stock.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewStockServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterStockServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//_ = consul.RegisterService(c.ListenOn, c.Consul) //if consul use

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
