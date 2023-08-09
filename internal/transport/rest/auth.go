package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/models"
	"net/http"
)

// @Summary User registration
// @Tags auth
// @Description signUp registers a new user
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ok"
// @Router /auth [post]
func (h *handler) signUp(ctx *gin.Context) {
	inputBodySignUp, _ := ctx.Get(inputSignUp)
	input, _ := inputBodySignUp.(models.SignUpInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration is not completed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sign up successful!"})
}

func (h *handler) signIn(ctx *gin.Context) {
	inputBodySignIn, _ := ctx.Get(inputSignIn)
	input, _ := inputBodySignIn.(models.SignInInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignInHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(ctx, input)
	if err != nil {
		logger.LogError(logger.SignInHandler, err)
		if errors.Is(err, service.UserNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "you're unauthorized"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "authentication is not completed"})
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *handler) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logger.LogError(logger.RefreshHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh-token not found"})
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx, cookie)
	if err != nil {
		logger.LogError(logger.RefreshHandler, err)
		if errors.Is(err, service.UserNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "you need to sign up"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create JWT"})
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}
