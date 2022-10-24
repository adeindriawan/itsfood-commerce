package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/utils"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func AuthorizedUser() gin.HandlerFunc {
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
		var user models.User
		if err := services.DB.Where("id = ?", userId).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": userId,
				"description": "Gagal mengambil data user dengan ID yang dimaksud.",
			})
			c.Abort()
			return
		}

		if user.Type != "Customer" {
			c.JSON(403, gin.H{
				"status": "failed",
				"errors": "User bukan merupakan Customer.",
				"result": nil,
				"description": "Tidak dapat melanjutkan request karena User bukan merupakan Customer.",
			})
			c.Abort()
			return
		} else {
			c.Set("user", user) // add user object to the context so it can be brought to next middleware
			c.Next()
		}
	}
}