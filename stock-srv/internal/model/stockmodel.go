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
	stockFieldNames          = builder.RawFieldNames(&Stock{})
	stockRows                = strings.Join(stockFieldNames, ",")
	stockRowsExpectAutoSet   = strings.Join(stringx.Remove(stockFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	stockRowsWithPlaceHolder = strings.Join(stringx.Remove(stockFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	StockModel interface {
		DecuctStock(goodsId , num int64) error
		AddStock(goodsId , num int64) error
		Insert(data *Stock) (sql.Result, error)
		FindOne(id int64) (*Stock, error)
		Update(data *Stock) error
		Delete(id int64) error
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


func (m *defaultStockModel) DecuctStock(goodsId , num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` - ? where `goods_id` = ? and num > 0", m.table)
	_, err := m.conn.Exec(query, num, goodsId)
	return err
}



func (m *defaultStockModel) AddStock(goodsId , num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` + ? where `goods_id` = ? and num > 0", m.table)
	_, err := m.conn.Exec(query, num, goodsId)
	return err
}

func (m *defaultStockModel) Insert(data *Stock) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, stockRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.GoodsId, data.Num)
	return ret, err
}

func (m *defaultStockModel) FindOne(id int64) (*Stock, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", stockRows, m.table)
	var resp Stock
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}


func (m *defaultStockModel) Update(data *Stock) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, stockRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.GoodsId, data.Num, data.Id)
	return err
}


func (m *defaultStockModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

/**
暴露给dtm barrier使用
*/
func (m *defaultStockModel) SqlDB()(*sql.DB, error) {
	return m.conn.RawDB()
}
