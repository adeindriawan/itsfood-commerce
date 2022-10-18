package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/adeindriawan/itsfood-commerce/controllers"
	"github.com/adeindriawan/itsfood-commerce/services"
)

func init() {
	services.InitRedis()
	services.InitMySQL()
}

func main() {
	r := gin.Default()
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	r.GET("/", func(c *gin.Context) {
		response := "This is ITSFood API Homepage. For full documentation, please visit this <a href='https://documenter.getpostman.com/view/2734100/2s83zdvRWQ' target='_blank'>link</a>"
		c.Data(200, "text/html; charset: utf-8", []byte(response))
	})
	r.POST("/todo", controllers.TokenAuthMiddleware(), controllers.CreateTodo)
	
	r.GET("/menus", controllers.GetMenus)
	r.GET("/menus/:id/details", controllers.GetMenuDetails)
	
	r.POST("/register", controllers.Register)
	r.POST("/admin/register", controllers.AdminRegister)
	r.POST("/customer/register", controllers.CustomerRegister)
	r.POST("/vendor/register", controllers.VendorRegister)
	r.POST("/customer/login", controllers.CustomerLogin)
	r.POST("/logout", controllers.TokenAuthMiddleware(), controllers.Logout)
	r.POST("/token/refresh", controllers.Refresh)
	r.POST("/password/forgot", controllers.ForgotPassword)
	r.POST("/password/reset", controllers.ResetPassword)

	r.POST("/cart", controllers.AddToCart)
	r.GET("/cart", controllers.ViewCart)
	r.PATCH("/cart", controllers.UpdateCart)
	r.DELETE("/cart", controllers.DeleteCart)
	r.GET("/cart/total", controllers.CartTotals)
	r.DELETE("/cart/destroy", controllers.DestroyCart)

	log.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
