package http

import (
	"net/http"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/adapter/jwt"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	usecases "github.com/AndreySmirnoffv/golang-auth-training/internal/usecases/users"
	"github.com/AndreySmirnoffv/golang-auth-training/pkg"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc usecases.UserUseCase
}

func NewUserHandler(uc usecases.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
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

	user := &entity.User{Email: req.Email, Password: req.Password}

	err := h.uc.Register(user)
	if err != nil {
		if err == usecases.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := userResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	resp := struct {
		ID           uint   `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, resp)
}

func JWTMiddleware(jwtSrv jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := pkg.ExtractBearerToken(c.GetHeader("Authorization"))
		refreshToken := c.GetHeader("X-Refresh-Token")

		claims, err := jwtSrv.ValidateAccessToken(accessToken)
		if err != nil && refreshToken != "" {
			newAccess, newRefresh, err := jwtSrv.RefreshTokens(refreshToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			c.Header("X-New-Access-Token", newAccess)
			c.Header("X-New-Refresh-Token", newRefresh)

			claims, _ = jwtSrv.ValidateAccessToken(newAccess)
		}

		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}
