package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/models"
	"net/http"
)

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
