package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"crud/models"
	"crud/pkg/helper"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (f *OrderRepo) Create(ctx context.Context, order *models.CreateOrder) (string, error) {

	var (
		id            = uuid.New().String()
		payed         sql.NullFloat64
		queryGetPrice string
		query         string
	)

	queryGetPrice = `
		SELECT
			price
		FROM
			books
		WHERE book_id = $1
	`

	query = `
		INSERT INTO orders(
			order_id,
			user_id,
			book_id,
			payed,
			updated_at
		) VALUES ( $1, $2, $3, $4, now() )
	`

	err := f.db.QueryRow(ctx, queryGetPrice, order.Book_id).
		Scan(
			&payed,
		)

	if err != nil {
		return "", err
	}

	_, err = f.db.Exec(ctx, query,
		id,
		order.User_id,
		order.Book_id,
		payed,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (f *OrderRepo) GetByPKey(ctx context.Context, pkey *models.OrderPrimarKey) (*models.Order, error) {

	var (
		id         sql.NullString
		user_id    sql.NullString
		book_id    sql.NullString
		payed      sql.NullFloat64
		created_at sql.NullString
	)

	query := `
		SELECT
			order_id,
			user_id,
			book_id,
			payed,
			updated_at
		FROM
			orders
		WHERE order_id = $1
	`

	err := f.db.QueryRow(ctx, query, pkey.Id).
		Scan(
			&id,
			&user_id,
			&book_id,
			&payed,
			&created_at,
		)

	if err != nil {
		return nil, err
	}

	return &models.Order{
		Id:        id.String,
		User_id:   user_id.String,
		Book_id:   book_id.String,
		Payed:     payed.Float64,
		CreatedAt: created_at.String,
	}, nil
}

func (f *OrderRepo) GetList(ctx context.Context, req *models.GetListOrderRequest) (*models.GetListOrderResponse, error) {

	var (
		resp   = models.GetListOrderResponse{}
		offset = ""
		limit  = ""
	)

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	query := `
		SELECT 
			COUNT(*) OVER(),
			users.first_name || ' ' || users.last_name as fullname,
			books.title,
			orders.payed,
			orders.created_at,
			orders.updated_at
		FROM
			orders
		JOIN users ON orders.user_id = users.user_id
		JOIN books ON orders.book_id = books.book_id
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)

	for rows.Next() {

		var (
			fullname   sql.NullString
			title      sql.NullString
			payed      sql.NullFloat64
			created_at sql.NullString
			updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&fullname,
			&title,
			&payed,
			&created_at,
			&updated_at,
		)

		if err != nil {
			return nil, err
		}

		resp.Orders = append(resp.Orders, &models.OrderGroup{
			FullName:  fullname.String,
			Title:     title.String,
			Payed:     payed.Float64,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})

	}

	return &resp, err
}

func (f *OrderRepo) Update(ctx context.Context, req *models.UpdateOrder) (int64, error) {

	var (
		query         = ""
		queryGetPrice = ""
		params        map[string]interface{}
	)

	queryGetPrice = `
		SELECT
			price
		FROM
			books
		WHERE book_id = $1
	`

	query = `
		UPDATE
			orders
		SET
			user_id = :user_id,
			book_id = :book_id,
			payed = :payed,
			updated_at = now()
		WHERE order_id = :order_id
	`

	err := f.db.QueryRow(ctx, queryGetPrice, req.Book_id).
		Scan(
			&req.Payed,
		)

	fmt.Println(req.Payed)
	if err != nil {
		return 0, err
	}

	params = map[string]interface{}{
		"order_id": req.Id,
		"user_id":  req.User_id,
		"book_id":  req.Book_id,
		"payed":    req.Payed,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *OrderRepo) Delete(ctx context.Context, req *models.OrderPrimarKey) error {

	_, err := f.db.Exec(ctx, "DELETE FROM orders WHERE order_id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
