package main

import (
	"flag"
	"fmt"
	//"github.com/zeromicro/zero-contrib/zrpc/registry/consul" //if consul use
	"gozerodtm/order-srv/internal/config"
	"gozerodtm/order-srv/internal/server"
	"gozerodtm/order-srv/internal/svc"
	"gozerodtm/order-srv/pb"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewOrderServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterOrderServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//_ = consul.RegisterService(c.ListenOn, c.Consul) //if consul use

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
