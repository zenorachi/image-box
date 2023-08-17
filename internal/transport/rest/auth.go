package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/model"
)

// @Summary User registration
// @Tags auth
// @Description signUp registers a new user
// @Accept  json
// @Produce  json
// @Param input body model.SignUpInput true "sign up info"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /auth/sign-up [post]
func (h *handler) signUp(ctx *gin.Context) {
	inputBodySignUp, _ := ctx.Get(inputSignUp)
	input, _ := inputBodySignUp.(model.SignUpInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userService.SignUp(ctx.Request.Context(), input); err != nil {
		logger.LogError(logger.SignUpHandler, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration is not completed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sign up successful!"})
}

// @Summary User authentication
// @Tags auth
// @Description signIn returns JWT token
// @Accept  json
// @Produce  json
// @Param input body model.SignInInput true "sign in info"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /auth/sign-in [post]
func (h *handler) signIn(ctx *gin.Context) {
	inputBodySignIn, _ := ctx.Get(inputSignIn)
	input, _ := inputBodySignIn.(model.SignInInput)

	if err := input.Validate(); err != nil {
		logger.LogError(logger.SignInHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(ctx.Request.Context(), input)
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

// @Summary Refresh JWT
// @Tags auth
// @Description refresh returns a new JWT token
// @Accept  json
// @Produce  json
// @HeaderParam Set-Cookie string true "RefreshToken"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /auth/refresh [get]
func (h *handler) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logger.LogError(logger.RefreshHandler, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh-token not found"})
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx.Request.Context(), cookie)
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
