package db

import (
	"context"
	"fmt"

	pb "github.com/opentelemetry/opentelemetry-demo/src/product-catalog/genproto/oteldemo"
)

type ProductRepository struct {
	db *DB
}

func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]*pb.Product, error) {
	query := `
		SELECT id, name, description, picture, 
		       price_units, price_nanos, price_currency_code, categories
		FROM products`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		var priceUnits int64
		var priceNanos int32
		var priceCurrencyCode string

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Picture,
			&priceUnits,
			&priceNanos,
			&priceCurrencyCode,
			&product.Categories,
		); err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}

		product.PriceUsd = &pb.Money{
			CurrencyCode: priceCurrencyCode,
			Units:        priceUnits,
			Nanos:        priceNanos,
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id string) (*pb.Product, error) {
	query := `
		SELECT id, name, description, picture, 
		       price_units, price_nanos, price_currency_code, categories
		FROM products
		WHERE id = $1`

	var product pb.Product
	var priceUnits int64
	var priceNanos int32
	var priceCurrencyCode string

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Picture,
		&priceUnits,
		&priceNanos,
		&priceCurrencyCode,
		&product.Categories,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	product.PriceUsd = &pb.Money{
		CurrencyCode: priceCurrencyCode,
		Units:        priceUnits,
		Nanos:        priceNanos,
	}

	return &product, nil
}

func (r *ProductRepository) SearchProducts(ctx context.Context, query string) ([]*pb.Product, error) {
	dbQuery := `
		SELECT id, name, description, picture, 
		       price_units, price_nanos, price_currency_code, categories
		FROM products
		WHERE name ILIKE $1 OR description ILIKE $1`

	searchPattern := "%" + query + "%"
	rows, err := r.db.Pool.Query(ctx, dbQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("error searching products: %w", err)
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		var priceUnits int64
		var priceNanos int32
		var priceCurrencyCode string

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Picture,
			&priceUnits,
			&priceNanos,
			&priceCurrencyCode,
			&product.Categories,
		); err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}

		product.PriceUsd = &pb.Money{
			CurrencyCode: priceCurrencyCode,
			Units:        priceUnits,
			Nanos:        priceNanos,
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
} 