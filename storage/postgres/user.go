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

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (f *UserRepo) Create(ctx context.Context, user *models.CreateUser) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO users(
			user_id,
			first_name,
			last_name,
			login,
			password,
			phone_number,
			balance,
			updated_at
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, now() )
	`

	_, err := f.db.Exec(ctx, query,
		id,
		user.First_name,
		user.Last_name,
		user.Login,
		user.Password,
		user.Phone_number,
		user.Balance,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (f *UserRepo) GetByPKey(ctx context.Context, pkey *models.UserPrimarKey) (*models.User, error) {

	var (
		id           sql.NullString
		first_name   sql.NullString
		last_name    sql.NullString
		login        sql.NullString
		password     sql.NullString
		phone_number sql.NullString
		balance      sql.NullFloat64
		createdAt    sql.NullString
		updatedAt    sql.NullString
	)

	if len(pkey.Login) > 0 {

		err := f.db.QueryRow(ctx, "SELECT user_id FROM users WHERE login = $1", pkey.Login).
			Scan(&pkey.Id)

		if err != nil {
			return nil, err
		}

	}

	query := `
		SELECT
			user_id,
			first_name,
			last_name,
			login,
			password,
			phone_number,
			balance,
			created_at,
			updated_at
		FROM
			users
		WHERE user_id = $1
	`

	err := f.db.QueryRow(ctx, query, pkey.Id).
		Scan(
			&id,
			&first_name,
			&last_name,
			&login,
			&password,
			&phone_number,
			&balance,
			&createdAt,
			&updatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:           id.String,
		First_name:   first_name.String,
		Last_name:    last_name.String,
		Login:        login.String,
		Password:     password.String,
		Phone_number: phone_number.String,
		Balance:      balance.Float64,
		CreatedAt:    createdAt.String,
		UpdatedAt:    updatedAt.String,
	}, nil
}

func (f *UserRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (*models.GetListUserResponse, error) {

	var (
		resp   = models.GetListUserResponse{}
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
			user_id,
			first_name,
			last_name,
			login,
			password,
			phone_number,
			balance,
			created_at,
			updated_at
		FROM
			users
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)

	for rows.Next() {

		var (
			id           sql.NullString
			first_name   sql.NullString
			last_name    sql.NullString
			login        sql.NullString
			password     sql.NullString
			phone_number sql.NullString
			balance      sql.NullFloat64
			createdAt    sql.NullString
			updatedAt    sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&last_name,
			&login,
			&password,
			&phone_number,
			&balance,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &models.User{
			Id:           id.String,
			First_name:   first_name.String,
			Last_name:    last_name.String,
			Login:        login.String,
			Password:     password.String,
			Phone_number: phone_number.String,
			Balance:      balance.Float64,
			CreatedAt:    createdAt.String,
			UpdatedAt:    updatedAt.String,
		})

	}

	return &resp, err
}

func (f *UserRepo) Update(ctx context.Context, req *models.UpdateUser) (int64, error) {

	var (
		query  = ""
		params map[string]interface{}
	)

	query = `
		UPDATE
			users
		SET
			first_name = :first_name,
			last_name = :last_name,
			login = :login,
			password = :password,
			phone_number = :phone_number,
			balance = :balance,
			updated_at = now()
		WHERE user_id = :user_id
	`

	params = map[string]interface{}{
		"user_id":      req.Id,
		"first_name":   req.First_name,
		"last_name":    req.Last_name,
		"login":        req.Login,
		"password":     req.Password,
		"phone_number": req.Phone_number,
		"balance":      req.Balance,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *UserRepo) Delete(ctx context.Context, req *models.UserPrimarKey) error {

	_, err := f.db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
