package handler

import (
	"crud/storage"
)

type HandlerV1 struct {
	storage storage.StorageI
}

func NewHandlerV1(storage storage.StorageI) *HandlerV1 {
	return &HandlerV1{
		storage: storage,
	}
}
