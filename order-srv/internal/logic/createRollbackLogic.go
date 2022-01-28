package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"gozerodtm/order-srv/internal/model"

	"github.com/yedf/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	fmt.Printf("订单回滚  , in: %+v \n", in)

	order, err := l.svcCtx.OrderModel.FindLastOneByUserIdGoodsId(in.UserId, in.GoodsId)
	if err != nil && err != model.ErrNotFound {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil, status.Error(codes.Internal, err.Error())
	}

	if order != nil {

		barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
		db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
		if err != nil {
			//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
			return nil, status.Error(codes.Internal, err.Error())
		}
		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {

			order.RowState = -1
			if err := l.svcCtx.OrderModel.Update(tx, order); err != nil {
				return fmt.Errorf("回滚订单失败  err : %v , userId:%d , goodsId:%d", err, in.UserId, in.GoodsId)
			}

			return nil
		}); err != nil {
			logx.Errorf("err : %v \n", err)

			//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
			return nil, status.Error(codes.Internal, err.Error())
		}

	}

	return &pb.CreateResp{}, nil
}
