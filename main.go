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
		response := "This is Itsfood Commerce Service API Homepage. For full documentation, please visit this <a href='https://documenter.getpostman.com/view/2734100/2s83zdvRWQ' target='_blank'>link</a>"
		c.Data(200, "text/html; charset: utf-8", []byte(response))
	})

	authorized := r.Group("/")
	authorized.Use(middlewares.Authorized())
	{
		authorized.POST("/todo", controllers.CreateTodo)
		authorized.POST("/logout", controllers.Logout)

		authorizedCustomer := authorized.Group("/")
		authorizedCustomer.Use(middlewares.AuthorizedCustomer())
		{
			authorizedCustomer.POST("/cart", controllers.AddToCart)
			authorizedCustomer.GET("/cart", controllers.ViewCart)
			authorizedCustomer.PATCH("/cart", controllers.UpdateCart)
			authorizedCustomer.DELETE("/cart", controllers.DeleteCart)
			authorizedCustomer.GET("/cart/total", controllers.CartTotals)
			authorizedCustomer.DELETE("/cart/destroy", controllers.DestroyCart)
			authorizedCustomer.GET("/orders/:id/details", controllers.OrderDetails)

			authorizedActiveCustomer := authorizedCustomer.Group("/")
			authorizedActiveCustomer.Use(middlewares.AuthorizedActiveCustomer())
			{
				authorizedActiveCustomer.GET("/orders", controllers.GetOrders)
				authorizedActiveCustomer.POST("/orders", controllers.CreateOrder)
				authorizedActiveCustomer.PATCH("/order-details/:id/accepted", controllers.MarkMenuAsAccepted)
				authorizedActiveCustomer.POST("v1/orders", controllers.CreateOrderV1)
			}
		}
	}
	
	r.GET("/menus", controllers.GetMenus)
	r.GET("/menus/:id/details", controllers.GetMenuDetails)
	r.GET("/vendors/:id/info", controllers.GetVendorInfo)
	
	r.POST("/register", controllers.Register)
	r.POST("/admin/register", controllers.AdminRegister)
	r.POST("/auth/register", controllers.CustomerRegister)
	r.POST("/vendor/register", controllers.VendorRegister)
	r.POST("/auth/login", controllers.CustomerLogin)
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/token/refresh", controllers.Refresh)
	r.POST("/password/forgot", controllers.ForgotPassword)
	r.POST("/password/reset", controllers.ResetPassword)

	log.Fatal(r.Run())
}
