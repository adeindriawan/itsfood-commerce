package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/adeindriawan/itsfood-commerce/utils"
	jwt "github.com/golang-jwt/jwt/v4"
)

func TokenValid(r *http.Request) error {
	token, err := utils.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Token dari user tidak valid.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}