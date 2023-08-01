package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"io"
	"log"
	"net/http"
)

func checkBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println("checkBody middleware", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Set("requestBody", body)
		ctx.Next()
	}
}

func checkJSONSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.Get("requestBody")
		body := requestBody.([]byte)
		var input models.SignUpInput
		if err := json.Unmarshal(body, &input); err != nil {
			log.Println("signUp middleware", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Set("inputSignUp", input)
		ctx.Next()
	}
}
