package main

import (
	"context"
	"crud/api"
	"crud/config"
	"crud/storage/postgres"
	"crud/storage/redis"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	storage, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.CloseDB()

	cache, err := redis.NewRedis(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.CloseDB()

	api.SetUpApi(&cfg, r, storage, cache)

	log.Printf("Listening port %v...\n", cfg.HTTPPort)
	err = r.Run(cfg.HTTPPort)
	if err != nil {
		panic(err)
	}
}
