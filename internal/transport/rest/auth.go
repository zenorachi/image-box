package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"log"
	"net/http"
)

func (h *handler) signUp(ctx *gin.Context) {
	inputSignUp, exists := ctx.Get("inputSignUp")
	if !exists {
		log.Println("signUp handler: inputData not found in context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	input, ok := inputSignUp.(models.SignUpInput)
	if !ok {
		log.Println("signUp handler: inputData has invalid type")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := input.Validate(); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		log.Println("signUp handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sign up successful!"})
}

func (h *handler) signIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignIn": "ok"})
}
