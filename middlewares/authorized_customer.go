package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/utils"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func AuthorizedCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.AuthCheck(c)
		if err != nil {
			c.JSON(403, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"results": nil,
				"description": "Gagal mengecek user ID pada request ini.",
			})
			c.Abort()
			return
		}
		
		var customer models.Customer
			if err := services.DB.Preload("Unit").Preload("User").Where("user_id = ?", userId).First(&customer).Error; err != nil {
				c.JSON(404, gin.H{
					"status": "failed",
					"errors": err.Error(),
					"result": userId,
					"description": "Gagal mengambil data customer dengan user ID yang dimaksud.",
				})
				c.Abort()
				return
			}
			c.Set("customer", customer) // add customer object to the context so it can be brought to next middleware
			c.Next()
	}
}