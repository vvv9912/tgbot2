package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"tgbotv2/internal/model"
)

type OrdersPostgresStorage struct {
	db *sqlx.DB
}

func NewOrdersPostgresStorage(db *sqlx.DB) *OrdersPostgresStorage {
	return &OrdersPostgresStorage{db: db}
}

func (s *OrdersPostgresStorage) OrdersByArticle() {

}

// todo
func (s *OrdersPostgresStorage) AddOrders(ctx context.Context, order model.Orders) error {
	//conn, err := s.db.Connx(ctx)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO orders (id_user, status_order, pvz, orderr, created_at, read_at)
	    				VALUES ($1, $2, $3, $4,$5, $6)
	    				ON CONFLICT DO NOTHING;`,
		order.TgID,
		order.StatusOrder,
		order.Pvz,
		order.Order,
		order.CreatedAt,
		order.ReadAt,
	); err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

type dbOrders struct {
	ID          int    `db:"id"`
	IDUser      int    `db:"id_user"`
	StatusOrder int    `db:"status_order"`
	Pvz         string `db:"pvz"`
	Order       string `db:"order"`
}
