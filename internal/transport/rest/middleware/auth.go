package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/models"
	"io"
	"log"
	"net/http"
)

func (mw *middleware) CheckBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println("checkBody middleware", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
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
			rest.LogError(rest.AuthMiddleware, err)
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
			rest.LogError(rest.AuthMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		}

		ctx.Set(InputSignIn, input)
		ctx.Next()
	}
}
