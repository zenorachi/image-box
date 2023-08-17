package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/model"
)

func (h *handler) CheckBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		}

		ctx.Set(requestBody, body)
		ctx.Next()
	}
}

func (h *handler) CheckJSONSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get(requestBody)
		body := requestBody.([]byte)
		var input model.SignUpInput
		if err := json.Unmarshal(body, &input); err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		}

		ctx.Set(inputSignUp, input)
		ctx.Next()
	}
}

func (h *handler) CheckJSONSignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get(requestBody)
		body := requestBody.([]byte)
		var input model.SignInInput
		if err := json.Unmarshal(body, &input); err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		}

		ctx.Set(inputSignIn, input)
		ctx.Next()
	}
}

func (h *handler) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getTokenFromRequest(ctx)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		id, err := h.userService.ParseToken(ctx.Request.Context(), token)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
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
