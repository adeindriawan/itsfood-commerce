package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func AuthorizedActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet("user").(models.User)
		if u.Status != "Activated" {
			c.JSON(403, gin.H{
				"status": "failed",
				"errors": "User sedang berstatus tidak aktif.",
				"result": nil,
				"description": "Tidak dapat melanjutkan request karena User berstatus tidak aktif.",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}