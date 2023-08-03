package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/models"
	"log"
	"net/http"
)

func (h *handler) signUp(ctx *gin.Context) {
	inputBodySignUp, _ := ctx.Get(inputSignUp)

	input, _ := inputBodySignUp.(models.SignUpInput)

	if err := input.Validate(); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sign up successful!"})
}

func (h *handler) signIn(ctx *gin.Context) {
	inputBodySignIn, _ := ctx.Get(inputSignIn)

	input, _ := inputBodySignIn.(models.SignInInput)

	if err := input.Validate(); err != nil {
		log.Println("signIn handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.userService.SignIn(ctx, input)
	if err != nil {
		if errors.Is(err, service.UserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "You need to Sign up"})
		}
		log.Println("signIn handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sign in successful!",
		"token": token})
}
