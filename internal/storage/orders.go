package storage

import "github.com/jmoiron/sqlx"

type OrdersPostgresStorage struct {
	db *sqlx.DB
}

func (s *OrdersPostgresStorage) OrdersByArticle() {

}

type dbOrders struct {
	ID          int    `db:"id"`
	IDUser      int    `db:"id_user"`
	StatusOrder int    `db:"status_order"`
	Pvz         string `db:"pvz"`
	Order       string `db:"order"`
}
