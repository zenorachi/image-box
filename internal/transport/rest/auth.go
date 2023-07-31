package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *handler) signUp(ctx *gin.Context) {
	//TODO: MIDDLEWARE

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	var input models.SignUpInput
	if err = json.Unmarshal(body, &input); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = h.userService.SignUp(ctx, input); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sign up successful!"})
}

func (h *handler) signIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignIn": "ok"})
}
