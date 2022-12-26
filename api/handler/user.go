package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"crud/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// CreateUser godoc
// @ID create_user
// @Router /user [POST]
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "CreateUserRequestBody"
// @Success 201 {object} models.User "GetUserBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var user models.CreateUser

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.User().Create(context.Background(), &user)
	if err != nil {
		log.Printf("error whiling Create: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling Create").Error())
		return
	}

	resp, err := h.storage.User().GetByPKey(
		context.Background(),
		&models.UserPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	err = h.cache.User().Delete(context.Background())

	if err != nil {
		log.Printf("error whiling cache delete: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling cache delete").Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetByIdUser godoc
// @ID get_by_id_user
// @Router /user/{id} [GET]
// @Summary Get By Id User
// @Description Get By Id User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.User "GetUserBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetUserById(c *gin.Context) {

	id := c.Param("id")

	resp, err := h.storage.User().GetByPKey(
		context.Background(),
		&models.UserPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetListUser godoc
// @ID get_list_user
// @Router /user [GET]
// @Summary Get List User
// @Description Get List User
// @Tags User
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} models.GetListUserResponse "GetUserBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetUserList(c *gin.Context) {
	var (
		limit  int
		offset int
		err    error
	)

	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("error whiling limit: %v\n", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	offsetStr := c.Query("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("error whiling limit: %v\n", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	users, err := h.cache.User().GetList(context.Background())
	if err == redis.Nil {

		resp, err := h.storage.User().GetList(
			context.Background(),
			&models.GetListUserRequest{
				Limit:  int32(limit),
				Offset: int32(offset),
			},
		)

		if err != nil {
			log.Printf("error whiling get list: %v", err)
			c.JSON(http.StatusInternalServerError, errors.New("error whiling get list").Error())
			return
		}

		fmt.Println("POSTGRES")

		err = h.cache.User().Create(context.Background(), resp)

		if err != nil {
			log.Printf("error whiling create cache list: %v", err)
			c.JSON(http.StatusInternalServerError, errors.New("error whiling create cache list").Error())
			return
		}

		c.JSON(http.StatusOK, resp)
	} else {

		if err != nil {
			log.Printf("error whiling get list cache: %v", err)
			c.JSON(http.StatusInternalServerError, errors.New("error whiling get list cache").Error())
			return
		}

		fmt.Println("REDIS")

		c.JSON(http.StatusOK, users)
	}

}

// UpdateUser godoc
// @ID update_user
// @Router /user/{id} [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param user body models.UpdateUserSwagger true "CreateUserRequestBody"
// @Success 200 {object} models.User "GetUsersBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) UpdateUser(c *gin.Context) {

	var (
		user models.UpdateUser
	)

	id := c.Param("id")

	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required user id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required user id").Error())
		return
	}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("error whiling update: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.Id = id

	rowsAffected, err := h.storage.User().Update(
		context.Background(),
		&user,
	)

	if err != nil {
		log.Printf("error whiling update: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update").Error())
		return
	}

	if rowsAffected == 0 {
		log.Printf("error whiling update rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update rows affected").Error())
		return
	}

	resp, err := h.storage.User().GetByPKey(
		context.Background(),
		&models.UserPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	respList, err := h.storage.User().GetList(
		context.Background(),
		&models.GetListUserRequest{
			Limit:  int32(0),
			Offset: int32(0),
		},
	)
	if err != nil {
		log.Printf("error whiling get list: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling get list").Error())
		return
	}

	err = h.cache.User().Update(context.Background(), respList)
	if err != nil {
		log.Printf("error whiling update cache list: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update cache list").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteByIdUser godoc
// @ID delete_by_id_user
// @Router /user/{id} [DELETE]
// @Summary Delete By Id User
// @Description Delete By Id User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.User "GetUserBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) DeleteUser(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required user id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required user id").Error())
		return
	}

	err := h.storage.User().Delete(
		context.Background(),
		&models.UserPrimarKey{
			Id: id,
		},
	)

	if err != nil {
		log.Printf("error whiling delete: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling delete").Error())
		return
	}

	respList, err := h.storage.User().GetList(
		context.Background(),
		&models.GetListUserRequest{
			Limit:  int32(0),
			Offset: int32(0),
		},
	)
	if err != nil {
		log.Printf("error whiling get list: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling get list").Error())
		return
	}

	err = h.cache.User().Update(context.Background(), respList)
	if err != nil {
		log.Printf("error whiling update cache list: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update cache list").Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
