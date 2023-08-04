package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"log"
	"net/http"
)

func (mw *middleware) CheckJSONSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get(requestBody)
		body := requestBody.([]byte)
		var input models.SignUpInput
		if err := json.Unmarshal(body, &input); err != nil {
			log.Println("signUp middleware", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
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
			log.Println("signIn middleware", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Set(InputSignIn, input)
		ctx.Next()
	}
}
