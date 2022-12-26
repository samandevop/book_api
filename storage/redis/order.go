package redis

import (
	"context"
	"crud/models"
	"encoding/json"

	"github.com/go-redis/redis"
)

type OrderRepo struct {
	client *redis.Client
}

func NewOrderRepo(client *redis.Client) *OrderRepo {
	return &OrderRepo{
		client: client,
	}
}

func (u *OrderRepo) Create(ctx context.Context, req *models.GetListOrderResponse) error {

	orders, err := json.Marshal(req.Orders)
	if err != nil {
		return err
	}

	err = u.client.Set("orders", orders, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderRepo) GetList(ctx context.Context) (*models.GetListOrderResponse, error) {

	var resp = models.GetListOrderResponse{}

	orders, err := u.client.Get("orders").Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(orders), &resp.Orders)
	if err != nil {
		return nil, err
	}

	resp.Count = int32(len(resp.Orders))

	return &resp, nil
}

func (u *OrderRepo) Update(ctx context.Context, req *models.GetListOrderResponse) error {

	orders, err := json.Marshal(req.Orders)
	if err != nil {
		return err
	}

	err = u.client.Set("orders", orders, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderRepo) Delete(ctx context.Context) error {

	err := u.client.Del("orders").Err()
	return err
}
