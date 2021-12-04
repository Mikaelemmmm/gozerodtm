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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *pb.CreateReq) (*pb.CreateResp, error) {

	fmt.Printf("创建订单 in : %+v \n" , in)

	//barrier防止空补偿、空悬挂等具体看dtm官网即可，别忘记加barrier表在当前库中，因为判断补偿与要执行的sql一起本地事务
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

		order := new(model.Order)
		order.GoodsId = in.GoodsId
		order.Num = in.Num
		order.UserId = in.UserId

		_, err = l.svcCtx.OrderModel.Insert(order)
		if err != nil {
			return fmt.Errorf("创建订单失败 err : %v , order:%+v \n", err, order)
		}

		return nil
	});err != nil{
		logx.Errorf("err : %v \n" , err)

		//这里如果想dtm回滚 ，grpc返回错误必须是codes.Aborted,dtmcli.ResultFailure，dtm就是根据这个判断的，没有为什么
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &pb.CreateResp{}, nil
}
