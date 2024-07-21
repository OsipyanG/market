package handler

import (
	"context"
	"net/http"

	auth "github.com/OsipyanG/market/protos/auth"
	jwt "github.com/OsipyanG/market/protos/jwt"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authClient      auth.AuthClient
	adminAuthClient auth.AuthAdminClient
}

func NewAuthHandler(a auth.AuthClient, ad auth.AuthAdminClient) *AuthHandler {
	return &AuthHandler{authClient: a, adminAuthClient: ad}
}

func (h *AuthHandler) NewUser(ctx *gin.Context) {
	userCredentials := &auth.UserCredentials{}

	if err := ctx.ShouldBindJSON(userCredentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	tokens, err := h.authClient.NewUser(context.TODO(), userCredentials)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	userCredentials := &auth.UserCredentials{}

	if err := ctx.ShouldBindJSON(userCredentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	tokens, err := h.authClient.Login(context.TODO(), userCredentials)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) UpdateTokens(ctx *gin.Context) {
	refreshToken := &auth.RefreshToken{}

	if err := ctx.ShouldBindJSON(refreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	tokens, err := h.authClient.UpdateTokens(context.TODO(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) UpdatePassword(ctx *gin.Context) {
	requestUpdatePassword := &auth.RequestUpdatePassword{}

	if err := ctx.ShouldBindJSON(requestUpdatePassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := ctx.MustGet("claims")

	jwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestUpdatePassword.JwtClaims = jwtClaims

	_, err := h.authClient.UpdatePassword(context.TODO(), requestUpdatePassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, "password updated")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	refreshToken := &auth.RefreshToken{}

	if err := ctx.ShouldBindJSON(refreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	_, err := h.authClient.Logout(context.TODO(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, "logout")
}

func (h *AuthHandler) GetAllUsersWithLevel(ctx *gin.Context) {
	requestByLevel := &auth.RequestByLevel{}

	if err := ctx.ShouldBindJSON(requestByLevel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := ctx.MustGet("claims")

	jwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestByLevel.JwtClaims = jwtClaims

	users, err := h.adminAuthClient.GetAllUsersWithLevel(context.TODO(), requestByLevel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *AuthHandler) SetAccessLevel(ctx *gin.Context) {
	requestSetAccessLevel := &auth.SetAccessLevelRequest{}

	if err := ctx.ShouldBindJSON(requestSetAccessLevel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	userID := ctx.Param("user_id")

	claims := ctx.MustGet("claims")

	jwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestSetAccessLevel.JwtClaims = jwtClaims
	requestSetAccessLevel.UserId = userID

	_, err := h.adminAuthClient.SetAccessLevel(context.TODO(), requestSetAccessLevel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, "access level updated")
}

func (h *AuthHandler) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})

		return
	}

	claims := ctx.MustGet("claims")

	jwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestDeleteUser := &auth.RequestByUserID{
		JwtClaims: jwtClaims,
		UserId:    userID,
	}

	_, err := h.adminAuthClient.DeleteUser(context.TODO(), requestDeleteUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, "user deleted")
}
