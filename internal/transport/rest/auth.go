package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest/middleware"
	"github.com/zenorachi/image-box/models"
	"log"
	"net/http"
	"strings"
)

func (h *handler) signUp(ctx *gin.Context) {
	inputBodySignUp, _ := ctx.Get(middleware.InputSignUp)

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
	inputBodySignIn, _ := ctx.Get(middleware.InputSignIn)

	input, _ := inputBodySignIn.(models.SignInInput)

	if err := input.Validate(); err != nil {
		log.Println("signIn handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(ctx, input)
	if err != nil {
		if errors.Is(err, service.UserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "You need to Sign up"})
		}
		log.Println("signIn handler", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"message": "Sign in successful!",
		"token": accessToken})
}

func (h *handler) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		log.Println("refresh cookie error")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx, cookie)
	if err != nil {
		log.Println("refresh service error")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *handler) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getTokenFromRequest(ctx)
		if err != nil {
			log.Println("auth middleware")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error"})
			return
		}

		id, err := h.userService.ParseToken(ctx, token)
		if err != nil {
			log.Println("auth middleware", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Set("userID", id)
		fmt.Println(id)
		ctx.Next()
	}
}

func getTokenFromRequest(ctx *gin.Context) (string, error) {
	header := ctx.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
