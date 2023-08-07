package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/models"
	"io"
	"net/http"
)

func (mw *middleware) CheckBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		ctx.Set(requestBody, body)
		ctx.Next()
	}
}

func (mw *middleware) CheckJSONSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get(requestBody)
		body := requestBody.([]byte)
		var input models.SignUpInput
		if err := json.Unmarshal(body, &input); err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		}

		ctx.Set(InputSignUp, input)
		ctx.Next()
	}
}

func (mw *middleware) CheckJSONSignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get(requestBody)
		body := requestBody.([]byte)
		var input models.SignInInput
		if err := json.Unmarshal(body, &input); err != nil {
			logger.LogError(logger.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		}

		ctx.Set(InputSignIn, input)
		ctx.Next()
	}
}
