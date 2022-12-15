package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/helper"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA)+1:]
		token, err := helper.VerifyJWT(tokenString)
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// token invalid
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				return
			case jwt.ValidationErrorExpired:
				// token expired
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized, token expired",
				})
				return
			default:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				return
			}
		}

		if token.Valid {
			claims := token.Claims.(*helper.JWTClaim)

			c.Set("UserID", claims.UserID)
			c.Set("UserEmail", claims.UserEmail)
			c.Set("UserRole", claims.UserRole)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthorizeUserRole(role string) gin.HandlerFunc {
	// use this only after AuthorizeJWT
	return func(c *gin.Context) {
		roleJWT := c.GetString("UserRole")
		fmt.Println("Role :", role, "--- roleJWT :", roleJWT)
		if role != roleJWT {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
