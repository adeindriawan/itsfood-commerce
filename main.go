package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/controllers"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/middlewares"
	"github.com/adeindriawan/itsfood-commerce/utils"
)

func init() {
	utils.LoadEnvVars()
	services.InitRedis()
	services.InitMySQL()
}

func main() {
	r := gin.Default()
	r.Use(utils.UseCORS())

	r.GET("/", func(c *gin.Context) {
		response := "This is ITSFood API Homepage. For full documentation, please visit this <a href='https://documenter.getpostman.com/view/2734100/2s83zdvRWQ' target='_blank'>link</a>"
		c.Data(200, "text/html; charset: utf-8", []byte(response))
	})

	authorized := r.Group("/")
	authorized.Use(middlewares.Authorized())
	{
		authorized.POST("/todo", controllers.CreateTodo)
		authorized.POST("/logout", controllers.Logout)

		authorizedUser := authorized.Group("/")
		authorizedUser.Use(middlewares.AuthorizedCustomer())
		{
			authorizedUser.POST("/cart", controllers.AddToCart)
			authorizedUser.GET("/cart", controllers.ViewCart)
			authorizedUser.PATCH("/cart", controllers.UpdateCart)
			authorizedUser.DELETE("/cart", controllers.DeleteCart)
			authorizedUser.GET("/cart/total", controllers.CartTotals)
			authorizedUser.DELETE("/cart/destroy", controllers.DestroyCart)
			authorizedUser.GET("/orders/:id/details", controllers.OrderDetails)

			authorizedActiveUser := authorizedUser.Group("/")
			authorizedActiveUser.Use(middlewares.AuthorizedActiveUser())
			{
				authorizedActiveUser.POST("/orders", controllers.CreateOrder)
			}
		}
	}
	
	r.GET("/menus", controllers.GetMenus)
	r.GET("/menus/:id/details", controllers.GetMenuDetails)
	
	r.POST("/register", controllers.Register)
	r.POST("/admin/register", controllers.AdminRegister)
	r.POST("/customer/register", controllers.CustomerRegister)
	r.POST("/vendor/register", controllers.VendorRegister)
	r.POST("/customer/login", controllers.CustomerLogin)
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/token/refresh", controllers.Refresh)
	r.POST("/password/forgot", controllers.ForgotPassword)
	r.POST("/password/reset", controllers.ResetPassword)

	log.Fatal(r.Run())
}
