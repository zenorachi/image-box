package middleware

import (
	"github.com/gin-gonic/gin"
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
