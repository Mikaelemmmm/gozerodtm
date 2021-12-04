package logic

import (
	"context"
	"fmt"
	"github.com/yedf/dtmcli"
	"github.com/yedf/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gozerodtm/stock-srv/internal/model"
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

	stock, err := l.svcCtx.StockModel.FindOneByGoodsId(in.GoodsId)
	if err != nil && err != model.ErrNotFound{
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil,status.Error(codes.Internal,err.Error())
	}
	if stock == nil || stock.Num < in.Num {
		//【回滚】库存不足确定需要dtm直接回滚，直接返回 codes.Aborted, dtmcli.ResultFailure 才可以回滚
		return nil,status.Error(codes.Aborted,dtmcli.ResultFailure)
	}

	//barrier防止空补偿、空悬挂等具体看dtm官网即可，别忘记加barrier表在当前库中，因为判断补偿与要执行的sql一起本地事务
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := l.svcCtx.StockModel.SqlDB()
	if err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil,status.Error(codes.Internal,err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil,status.Error(codes.Internal,err.Error())
	}
	if err := barrier.Call(tx, func(db dtmcli.DB) error {

		if err := l.svcCtx.StockModel.DecuctStock(tx,in.GoodsId, in.Num);err!= nil{
			return fmt.Errorf("扣库存失败 err : %v , in:%+v \n",err,in)
		}

		//！！开启测试！！ ： 测试订单回滚更改状态为失效，并且当前库扣失败不需要回滚
		//return fmt.Errorf("扣库存失败 err : %v , in:%+v \n",err,in)

		return nil
	});err != nil{
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil,status.Error(codes.Internal,err.Error())
	}

	return &pb.DeductResp{}, nil
}
