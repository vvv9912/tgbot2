package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"tgbotv2/internal/model"
)

type ProductsPostgresStorage struct {
	db *sqlx.DB
}

func NewProductsPostgresStorage(db *sqlx.DB) *ProductsPostgresStorage {
	return &ProductsPostgresStorage{db: db}
}
func (s *ProductsPostgresStorage) AddProduct(ctx context.Context, product model.Products) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO products (article, catalog, name, description, photo_url, price, length, width, heigth, weight)
	    				VALUES ($1, $2, $3,$4, $5, $6,$7, $8, $9, $10)
	    				ON CONFLICT DO NOTHING;`,
		product.Article,
		product.Catalog,
		product.Name,
		product.Description,
		product.PhotoUrl,
		product.Price,
		product.Length,
		product.Width,
		product.Height,
		product.Weight,
		//users.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ProductsPostgresStorage) Catalog(ctx context.Context) ([]string, error) {

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var catalog []string
	if err := conn.SelectContext(ctx, &catalog, `SELECT DISTINCT catalog FROM products`); err != nil {
		return nil, err
	}

	return catalog, nil
}
func (s *ProductsPostgresStorage) ProductsByCatalog(ctx context.Context, ctlg string) ([]model.Products, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//var products []model.Products
	var products []dbProduct
	if err := conn.SelectContext(ctx,
		&products,
		`SELECT
     article AS a_article,
     catalog AS a_catalog,
     name AS a_name,
     description AS a_description,
     photo_url AS a_photo_url,
     price AS a_price
	 FROM products
	 WHERE catalog = $1`,
		ctlg); err != nil {
		return nil, err
	}

	return lo.Map(products, func(product dbProduct, _ int) model.Products { return model.Products(product) }), nil
}
func (s *ProductsPostgresStorage) ProductByArticle(ctx context.Context, article int) (model.Products, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return model.Products{}, err
	}
	defer conn.Close()
	//var products []model.Products
	var products dbProduct
	row := conn.QueryRowContext(ctx,
		`SELECT
     			article AS c_article,
     			catalog AS c_catalog,
     			name AS c_name,
     			description AS c_description,
     			photo_url AS c_photo_url,
     			price AS c_price,
     			length AS c_length,
     			width AS c_width,
     			heigth AS c_height,
     			weight AS c_weight
	 			FROM products
	 			WHERE (article = $1)`,
		article)
	err = row.Scan(
		&products.Article,
		&products.Catalog,
		&products.Name,
		&products.Description,
		&products.PhotoUrl,
		&products.Price,
		&products.Length,
		&products.Width,
		&products.Height,
		&products.Weight)
	if err != nil {
		return model.Products{}, err
	}

	return model.Products{
		Article:     products.Article,
		Catalog:     products.Catalog,
		Name:        products.Name,
		Description: products.Description,
		PhotoUrl:    products.PhotoUrl,
		Price:       products.Price,
		Length:      products.Length,
		Width:       products.Width,
		Height:      products.Height,
		Weight:      products.Weight,
	}, err
}

type dbProduct struct {
	Article     int     `db:"a_article"`
	Catalog     string  `db:"a_catalog"`
	Name        string  `db:"a_name"`
	Description string  `db:"a_description"`
	PhotoUrl    []byte  `db:"a_photo_url"`
	Price       float64 `db:"a_price"`
	Length      int     `db:"a_length"`
	Width       int     `db:"a_width"`
	Height      int     `db:"a_height"`
	Weight      int     `db:"a_weight"`
}
