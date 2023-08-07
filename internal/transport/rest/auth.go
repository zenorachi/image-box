package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/internal/transport/rest/middleware"
	"github.com/zenorachi/image-box/models"
	"net/http"
	"strings"
)

func (h *handler) signUp(ctx *gin.Context) {
	inputBodySignUp, _ := ctx.Get(middleware.InputSignUp)
	input, _ := inputBodySignUp.(models.SignUpInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration is not completed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sign up successful!"})
}

func (h *handler) signIn(ctx *gin.Context) {
	inputBodySignIn, _ := ctx.Get(middleware.InputSignIn)
	input, _ := inputBodySignIn.(models.SignInInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignInHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(ctx, input)
	if err != nil {
		logger.LogError(logger.SignInHandler, err)
		if errors.Is(err, service.UserNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "you're unauthorized"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "authentication is not completed"})
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *handler) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logger.LogError(logger.RefreshHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh-token not found"})
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx, cookie)
	if err != nil {
		logger.LogError(logger.RefreshHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot create JWT"})
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *handler) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getTokenFromRequest(ctx)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		id, err := h.userService.ParseToken(ctx, token)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot parse token"})
		}

		ctx.Set("userID", id)
		ctx.Next()
	}
}

func getTokenFromRequest(ctx *gin.Context) (string, error) {
	header := ctx.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
