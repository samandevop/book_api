package handler

import (
	"crud/config"
	"crud/storage"
)

type HandlerV1 struct {
	cfg     *config.Config
	storage storage.StorageI
	cache   storage.CacheI
}

func NewHandlerV1(cfg *config.Config, storage storage.StorageI, cache storage.CacheI) *HandlerV1 {
	return &HandlerV1{
		cfg:     cfg,
		storage: storage,
		cache:   cache,
	}
}
