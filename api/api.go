package api

import (
	_ "crud/api/docs"
	"crud/api/handler"
	"crud/config"
	"crud/pkg/helper"
	"crud/storage"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpApi(cfg *config.Config, r *gin.Engine, storage storage.StorageI, cache storage.CacheI) {

	handlerV1 := handler.NewHandlerV1(cfg, storage, cache)

	r.Use(customCORSMiddleware())

	// v1 := r.Group("/v1")
	// v2 := r.Group("/v2")

	r.POST("/login", handlerV1.Login)
	r.POST("/loginsuper", handlerV1.LoginSuper)

	r.POST("/refreshclienttoken")

	r.Use(checkTokenSuper())
	r.Use(checkTokenClient())
	r.POST("/book", handlerV1.CreateBook)
	r.GET("/book/:id", handlerV1.GetBookById)
	r.GET("/book", handlerV1.GetBookList)
	r.PUT("/book/:id", handlerV1.UpdateBook)
	r.DELETE("/book/:id", handlerV1.DeleteBook)

	r.POST("/user", handlerV1.CreateUser)
	r.GET("/user/:id", handlerV1.GetUserById)
	r.GET("/user", handlerV1.GetUserList)
	r.PUT("/user/:id", handlerV1.UpdateUser)
	r.DELETE("/user/:id", handlerV1.DeleteUser)

	r.POST("/order", handlerV1.CreateOrder)
	r.GET("/order/:id", handlerV1.GetOrderById)
	r.GET("/order", handlerV1.GetOrderList)
	r.PUT("/order/:id", handlerV1.UpdateOrder)
	r.DELETE("/order/:id", handlerV1.DeleteOrder)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func checkTokenSuper() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := ctx.Request.Header["Authorization"]; ok {
			_, err := helper.ExtractClaims(ctx.Request.Header["Authorization"][0], config.Load().AuthSecretKey)
			_, err2 := helper.ExtractClaims(ctx.Request.Header["Authorization"][0], config.Load().SuperAdmin)
			if err != nil && err2 != nil {
				ctx.AbortWithError(http.StatusForbidden, errors.New("not found password"))
				return
			} else {
				ctx.Next()
			}
		}
	}
}

func checkTokenClient() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := ctx.Request.Header["Authorization"]; ok {
			_, err := helper.ExtractClaims(ctx.Request.Header["Authorization"][0], config.Load().AuthSecretKey)
			_, err2 := helper.ExtractClaims(ctx.Request.Header["Authorization"][0], config.Load().Client)
			if err != nil && err2 != nil {
				ctx.AbortWithError(http.StatusForbidden, errors.New("not found password"))
				return
			} else {
				ctx.Next()
			}
		}
	}
}

func customCORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
