package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/habibmrizki/finalphase3/pkg"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Message: "Authorization header required",
				Status:  "Unauthorized",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Message: "Invalid Authorization format. Expected 'Bearer <token>'",
				Status:  "Unauthorized",
			})
			return
		}

		tokenString := parts[1]

		claims := &pkg.Claims{}
		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Message: "Invalid or expired token",
				Status:  "Unauthorized",
			})
			return
		}

		ctx.Set("user_id", claims.UserId)

		ctx.Next()
	}
}
