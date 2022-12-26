package redis

import (
	"context"
	"crud/models"
	"encoding/json"

	"github.com/go-redis/redis"
)

type UserRepo struct {
	client *redis.Client
}

func NewUserRepo(client *redis.Client) *UserRepo {
	return &UserRepo{
		client: client,
	}
}

func (u *UserRepo) Create(ctx context.Context, req *models.GetListUserResponse) error {

	users, err := json.Marshal(req.Users)
	if err != nil {
		return err
	}

	err = u.client.Set("users", users, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetList(ctx context.Context) (*models.GetListUserResponse, error) {

	var resp = models.GetListUserResponse{}

	users, err := u.client.Get("users").Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(users), &resp.Users)
	if err != nil {
		return nil, err
	}

	resp.Count = int32(len(resp.Users))

	return &resp, nil
}

func (u *UserRepo) Update(ctx context.Context, req *models.GetListUserResponse) error {

	users, err := json.Marshal(req.Users)
	if err != nil {
		return err
	}

	err = u.client.Set("users", users, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) Delete(ctx context.Context) error {

	err := u.client.Del("users").Err()
	return err
}
