package handler

import (
	"context"
	"crud/config"
	"crud/models"
	"crud/pkg/helper"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @ID login
// @Router /login [POST]
// @Summary Create Login
// @Description Create Login
// @Tags Login
// @Accept json
// @Produce json
// @Param Login body models.Login true "LoginRequestBody"
// @Success 201 {object} models.LoginResponse "GetLoginBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) Login(c *gin.Context) {
	var login models.Login

	err := c.ShouldBindJSON(&login)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.storage.User().GetByPKey(
		context.Background(),
		&models.UserPrimarKey{Login: login.Login},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	data := map[string]interface{}{
		"user_id": resp.Id,
	}

	token, err := helper.GenerateJWT(data, config.TimeExpiredAt, h.cfg.AuthSecretKey)
	if err != nil {
		log.Printf("error whiling GenerateJWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GenerateJWT").Error())
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{AccessToken: token})
}
