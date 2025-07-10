package http

import (
	"net/http"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/usecases"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *usecases.UserUseCase
}

func NewUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &entity.User{Email: req.Email, Password: req.Password}

	if err := h.uc.Register(u); err != nil {
		if err != usecases.ErrEmailExists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
