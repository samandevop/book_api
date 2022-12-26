package storage

import (
	"context"
	"crud/models"
)

type CacheI interface {
	CloseDB()
	User() UserCacheI
	Order() OrderCacheI
}

type UserCacheI interface {
	Create(ctx context.Context, req *models.GetListUserResponse) error
	GetList(ctx context.Context) (*models.GetListUserResponse, error)
	Update(ctx context.Context, req *models.GetListUserResponse) error
	Delete(ctx context.Context) error
}

type OrderCacheI interface {
	Create(ctx context.Context, req *models.GetListOrderResponse) error
	GetList(ctx context.Context) (*models.GetListOrderResponse, error)
	Update(ctx context.Context, req *models.GetListOrderResponse) error
	Delete(ctx context.Context) error
}
