package logic

import (
	"context"
	"fmt"
	"github.com/yedf/dtmcli/dtmimp"
	"gozerodtm/order-srv/order"
	"gozerodtm/stock-srv/stock"
	"net/http"

	"gozerodtm/order-api/internal/svc"
	"gozerodtm/order-api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/yedf/dtmgrpc"
)


// dtm已经通过配置，注册到下面这个地址，因此在dtmgrpc中使用该地址
var dtmServer = "etcd://localhost:2379/dtmservice"

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateLogic {
	return CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req types.QuickCreateReq,r *http.Request) (*types.QuickCreateResp, error) {

	orderRpcBusiServer, err := l.svcCtx.Config.OrderRpcConf.BuildTarget()
	if err != nil{
		return nil,fmt.Errorf("下单异常超时")
	}
	stockRpcBusiServer, err := l.svcCtx.Config.StockRpcConf.BuildTarget()
	if err != nil{
		return nil,fmt.Errorf("下单异常超时")
	}

	createOrderReq:= &order.CreateReq{UserId: req.UserId,GoodsId: req.GoodsId,Num: req.Num}
	deductReq:= &stock.DecuctReq{GoodsId: req.GoodsId,Num: req.Num}

	//这里只举了saga例子，tcc等其他例子基本没啥区别具体可以看dtm官网

	gid := dtmgrpc.MustGenGid(dtmServer)
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(orderRpcBusiServer+"/pb.order/create", orderRpcBusiServer+"/pb.order/createRollback", createOrderReq).
		Add(stockRpcBusiServer+"/pb.stock/deduct", stockRpcBusiServer+"/pb.stock/deductRollback", deductReq)

	err = saga.Submit()
	dtmimp.FatalIfError(err)
	if err != nil{
		return nil,fmt.Errorf("submit data to  dtm-server err  : %+v \n",err)
	}

	return &types.QuickCreateResp{}, nil
}
