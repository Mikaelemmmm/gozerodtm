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

type DeductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductRollbackLogic {
	return &DeductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductRollbackLogic) DeductRollback(in *pb.DecuctReq) (*pb.DeductResp, error) {

	fmt.Printf("库存回滚 in : %+v \n" , in)

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

		if err := l.svcCtx.StockModel.AddStock(tx,in.GoodsId, in.Num);err!= nil{
			return fmt.Errorf("回滚库存失败 err : %v ,goodsId:%d , num :%d", err,in.GoodsId,in.Num)
		}
		return nil
	});err != nil{
		logx.Errorf("err : %v \n" , err)
		return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
	}

	return &pb.DeductResp{}, nil
}
