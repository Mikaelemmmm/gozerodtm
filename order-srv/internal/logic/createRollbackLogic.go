package logic

import (
	"context"
	"fmt"
	"github.com/yedf/dtmcli"
	"github.com/yedf/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gozerodtm/order-srv/internal/model"

	"gozerodtm/order-srv/internal/svc"
	"gozerodtm/order-srv/pb"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRollbackLogic {
	return &CreateRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRollbackLogic) CreateRollback(in *pb.CreateReq) (*pb.CreateResp, error) {

	fmt.Printf("订单回滚  , in: %+v \n",in)

	order, err := l.svcCtx.OrderModel.FindLastOneByUserIdGoodsId(in.UserId, in.GoodsId)
	if err != nil && err != model.ErrNotFound{
		logx.Errorf("FindLastOneByUserIdGoodsId err : %v \n" , err)
		//这里如果想dtm回滚 ，grpc返回错误必须是codes.Aborted,dtmcli.ResultFailure，dtm就是根据这个判断的，没有为什么
		return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
	}

	if order != nil{

		barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
		db, err := l.svcCtx.OrderModel.SqlDB()
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

			order.RowState = -1
			if err := l.svcCtx.OrderModel.Update(order);err!= nil{
				return fmt.Errorf("回滚订单失败  err : %v , userId:%d , goodsId:%d",err,in.UserId,in.GoodsId)
			}

			return nil
		});err != nil{
			logx.Errorf("err : %v \n" , err)
			//这里如果想dtm回滚 ，grpc返回错误必须是codes.Aborted,dtmcli.ResultFailure，dtm就是根据这个判断的，没有为什么
			return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
		}

	}




	return &pb.CreateResp{}, nil
}
