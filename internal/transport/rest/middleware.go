package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"io"
	"log"
	"net/http"
)

func checkJSONBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println("signUp middleware", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var input models.SignUpInput
		if err = json.Unmarshal(body, &input); err != nil {
			log.Println("signUp middleware", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Set("inputData", &input)
		ctx.Next()
	}
}
