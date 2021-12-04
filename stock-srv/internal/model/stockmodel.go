package model

import (
	"database/sql"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)


type (
	StockModel interface {
		DecuctStock(tx *sql.Tx,goodsId , num int64) error
		AddStock(tx *sql.Tx,goodsId , num int64) error
		SqlDB()(*sql.DB, error)
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


func (m *defaultStockModel) DecuctStock(tx *sql.Tx,goodsId , num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` - ? where `goods_id` = ? and num > 0", m.table)
	_, err := sqlx.NewSessionFromTx(tx).Exec(query, num, goodsId)
	return err
}

func (m *defaultStockModel) AddStock(tx *sql.Tx,goodsId , num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` + ? where `goods_id` = ?", m.table)
	_, err :=sqlx.NewSessionFromTx(tx).Exec(query, num, goodsId)
	return err
}

/**
暴露给dtm barrier使用
*/
func (m *defaultStockModel) SqlDB()(*sql.DB, error) {
	return m.conn.RawDB()
}
