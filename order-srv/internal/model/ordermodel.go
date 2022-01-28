package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	orderFieldNames          = builder.RawFieldNames(&Order{})
	orderRows                = strings.Join(orderFieldNames, ",")
	orderRowsExpectAutoSet   = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	orderRowsWithPlaceHolder = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	OrderModel interface {
		FindLastOneByUserIdGoodsId(userId,goodsId int64) (*Order, error)
		Insert(tx *sql.Tx,data *Order) (sql.Result, error)
		Update(tx *sql.Tx,data *Order) error
	}

	defaultOrderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Order struct {
		Id       int64 `db:"id"`
		UserId   int64 `db:"user_id"`
		GoodsId  int64 `db:"goods_id"`  // 商品id
		Num      int64 `db:"num"`       // 下单数量
		RowState int64 `db:"row_state"` // -1:下单回滚失效 0:待支付
	}
)

func NewOrderModel(conn sqlx.SqlConn) OrderModel {
	return &defaultOrderModel{
		conn:  conn,
		table: "`order`",
	}
}

func (m *defaultOrderModel) FindLastOneByUserIdGoodsId(userId,goodsId int64) (*Order, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and goods_id =? order by id desc limit 1 ", orderRows, m.table)
	var resp Order
	err := m.conn.QueryRow(&resp, query, userId,goodsId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderModel) Insert(tx *sql.Tx,data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ? ,?)", m.table, orderRowsExpectAutoSet)
	return tx.Exec(query, data.UserId, data.GoodsId, data.Num,data.RowState)
}

func (m *defaultOrderModel) Update(tx *sql.Tx,data *Order) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
	_, err := tx.Exec(query, data.UserId, data.GoodsId, data.Num,data.RowState, data.Id)
	return err
}

