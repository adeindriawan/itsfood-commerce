package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/adeindriawan/itsfood-commerce/controllers"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func main() {
	r := gin.Default()
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	r.GET("/", controllers.SimpleRequest0)
	r.GET("/ping", controllers.SimpleRequest1)
	r.GET("/pong", controllers.SimpleRequest2)

	// This handler will match /user/john but will not match /user/ or /user
	r.GET("/user/:name", controllers.SimpleRequest3)
	
	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	r.GET("/user/:name/*action", controllers.SimpleRequest4)
	r.POST("/user/:name/*action", controllers.SimpleRequest5)

	models.ConnectDB()
	r.GET("/products", controllers.FindProducts)
	r.GET("/products/:id/details", controllers.ProductDetails)
	r.POST("/products", controllers.CreateProduct)
	r.PATCH("/products/:id", controllers.UpdateProduct)
	r.GET("/menus", controllers.GetMenus)
	r.GET("/menus/:id/details", controllers.GetMenuDetails)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	r.POST("/todo", controllers.CreateTodo)

	log.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
