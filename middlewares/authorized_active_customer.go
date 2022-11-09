package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func AuthorizedActiveCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet("user").(models.User)
		if u.Status != "Activated" {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": "User sedang berstatus tidak aktif.",
				"result": nil,
				"description": "Tidak dapat melanjutkan request karena User berstatus tidak aktif.",
			})
			c.Abort()
			return
		}

		cust := c.MustGet("customer").(models.Customer)
		if cust.Status != "Activated" {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": "Customer sedang berstatus tidak aktif.",
				"result": nil,
				"description": "Tidak dapat melanjutkan request karena Customer berstatus tidak aktif.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}