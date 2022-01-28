package model

import (
	"database/sql"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"strings"
)

var (
	stockFieldNames          = builder.RawFieldNames(&Stock{})
	stockRows                = strings.Join(stockFieldNames, ",")
)

type (
	StockModel interface {
		FindOneByGoodsId(goodsId int64) (*Stock, error)
		DecuctStock(tx *sql.Tx,goodsId , num int64) (sql.Result,error)
		AddStock(tx *sql.Tx,goodsId , num int64) error
	}

	defaultStockModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Stock struct {
		Id      int64 `db:"id"`
		GoodsId int64 `db:"goods_id"` // 商品id
		Num     int64 `db:"num"`      // 库存数量
	}
)

func NewStockModel(conn sqlx.SqlConn) StockModel {
	return &defaultStockModel{
		conn:  conn,
		table: "`stock`",
	}
}

func (m *defaultStockModel) FindOneByGoodsId(goodsId int64) (*Stock, error) {
	query := fmt.Sprintf("select %s from %s where `goods_id` = ? limit 1", stockRows, m.table)
	var resp Stock
	err := m.conn.QueryRow(&resp, query, goodsId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}


func (m *defaultStockModel) DecuctStock(tx *sql.Tx,goodsId , num int64) (sql.Result,error) {
	query := fmt.Sprintf("update %s set `num` = `num` - ? where `goods_id` = ? and num >= ?", m.table)
	return tx.Exec(query,num, goodsId,num)

}

func (m *defaultStockModel) AddStock(tx *sql.Tx,goodsId , num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` + ? where `goods_id` = ?", m.table)
	_, err :=tx.Exec(query, num, goodsId)
	return err
}

