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

type BookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) *BookRepo {
	return &BookRepo{
		db: db,
	}
}

func (f *BookRepo) Create(ctx context.Context, book *models.CreateBook) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO books(
			book_id,
			title,
			author,
			price,
			updated_at
		) VALUES ( $1, $2, $3, $4, now() )
	`

	_, err := f.db.Exec(ctx, query,
		id,
		book.Title,
		book.Author,
		book.Price,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (f *BookRepo) GetByPKey(ctx context.Context, pkey *models.BookPrimarKey) (*models.Book, error) {

	var (
		id        sql.NullString
		title     sql.NullString
		author    sql.NullString
		price     sql.NullFloat64
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT
			book_id,
			title,
			author,
			price,
			created_at,
			updated_at
		FROM
			books
		WHERE book_id = $1
	`

	err := f.db.QueryRow(ctx, query, pkey.Id).
		Scan(
			&id,
			&title,
			&author,
			&price,
			&createdAt,
			&updatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &models.Book{
		Id:        id.String,
		Title:     title.String,
		Author:    author.String,
		Price:     price.Float64,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (f *BookRepo) GetList(ctx context.Context, req *models.GetListBookRequest) (*models.GetListBookResponse, error) {

	var (
		resp   = models.GetListBookResponse{}
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
			book_id,
			title,
			author,
			price,
			created_at,
			updated_at
		FROM
			books
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)

	for rows.Next() {

		var (
			id        sql.NullString
			title     sql.NullString
			author    sql.NullString
			price     sql.NullFloat64
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&title,
			&author,
			&price,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Books = append(resp.Books, &models.Book{
			Id:        id.String,
			Title:     title.String,
			Author:    author.String,
			Price:     price.Float64,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})

	}

	return &resp, err
}

func (f *BookRepo) Update(ctx context.Context, req *models.UpdateBook) (int64, error) {

	var (
		query  = ""
		params map[string]interface{}
	)

	query = `
		UPDATE
			books
		SET
			title = :title,
			author = :author,
			price = :price,
			updated_at = now()
		WHERE book_id = :book_id
	`

	params = map[string]interface{}{
		"book_id": req.Id,
		"title":   req.Title,
		"author":  req.Author,
		"price":   req.Price,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *BookRepo) Delete(ctx context.Context, req *models.BookPrimarKey) error {

	_, err := f.db.Exec(ctx, "DELETE FROM books WHERE book_id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
