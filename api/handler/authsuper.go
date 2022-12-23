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

// LoginSuper godoc
// @ID loginSuper
// @Router /loginsuper [POST]
// @Summary Create LoginSuper
// @Description Create LoginSuper
// @Tags LoginSuper
// @Accept json
// @Produce json
// @Param Login body models.Login true "LoginSuperRequestBody"
// @Success 201 {object} models.LoginResponse "GetLoginSuperBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) LoginSuper(c *gin.Context) {
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

	if login.Password != resp.Password {
		c.JSON(http.StatusInternalServerError, errors.New("error password is not correct").Error())
		return
	}

	data := map[string]interface{}{
		"user_id": resp.Id,
	}

	token, err := helper.GenerateJWT(data, config.SuperTimeExpiredAt, h.cfg.AuthSecretKey, h.cfg.SuperAdmin)
	if err != nil {
		log.Printf("error whiling GenerateJWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GenerateJWT").Error())
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{AccessToken: token})
}
