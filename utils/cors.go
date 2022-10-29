package utils

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func UseCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
				return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})
}