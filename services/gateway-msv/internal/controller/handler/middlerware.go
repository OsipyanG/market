package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/OsipyanG/market/protos/auth"
	"github.com/gin-gonic/gin"
)

const (
	AuthHeaderParts = 2
)

func AuthMiddleware(authClient auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not provided"})

			return
		}

		unpackedToken := strings.Split(token, " ")
		if len(unpackedToken) < AuthHeaderParts {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no bearer"})

			return
		}

		token = unpackedToken[1]

		jwtClaimsRequest := &auth.AccessToken{
			Value: token,
		}

		jwtClaims, err := authClient.GetJWTClaims(context.TODO(), jwtClaimsRequest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token is invalid:"})

			return
		}

		c.Set("token", token)
		c.Set("claims", jwtClaims)

		c.Next()
	}
}
