package api

import (
	_ "crud/api/docs"
	"crud/api/handler"
	"crud/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpApi(r *gin.Engine, storage storage.StorageI) {

	handlerV1 := handler.NewHandlerV1(storage)

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
