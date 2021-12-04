package logic

import (
	"context"
	"fmt"
	"github.com/yedf/dtmcli"
	"github.com/yedf/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gozerodtm/stock-srv/internal/svc"
	"gozerodtm/stock-srv/pb"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductLogic {
	return &DeductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductLogic) Deduct(in *pb.DecuctReq) (*pb.DeductResp, error) {

	fmt.Printf("扣库存start....")
	//barrier防止空补偿、空悬挂等具体看dtm官网即可，别忘记加barrier表在当前库中，因为判断补偿与要执行的sql一起本地事务
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := l.svcCtx.StockModel.SqlDB()
	if err != nil {
		logx.Errorf("获取sqlDB失败 err : %v",err)
		return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
	}
	tx, err := db.Begin()
	if err != nil {
		logx.Errorf("开启事务失败 err : %v",err)
		return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
	}
	if err := barrier.Call(tx, func(db dtmcli.DB) error {

		if err := l.svcCtx.StockModel.DecuctStock(tx,in.GoodsId, in.Num);err!= nil{
			return fmt.Errorf("扣库存失败 err : %v , in:%+v \n",err,in)
		}

		//！！开启测试！！ ： 测试订单回滚更改状态为失效，并且当前库扣失败不需要回滚
		//return fmt.Errorf("扣库存失败 err : %v , in:%+v \n",err,in)

		return nil
	});err != nil{
		logx.Errorf("err : %v \n" , err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &pb.DeductResp{}, nil
}
