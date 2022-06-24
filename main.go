package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/pong", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})

	// This handler will match /user/john but will not match /user/ or /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	r.POST("/user/:name/*action", func(c *gin.Context) {
		b := c.FullPath() == "/user/:name/*action" // true
		c.String(http.StatusOK, "%t", b)
	})

	r.GET("/test/route", testRoute)
	r.GET("/test/db/migrate", testDbMigrate)
	r.GET("/test/db/create", testDbCreate)
	r.GET("/test/db/product/:id/details", getProductDetails)
	r.GET("/test/db/products", getProducts)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func testRoute(c *gin.Context) {
	log.Println("separate function")
}

func testDbMigrate(c *gin.Context) {
	dsn := "host=localhost user=postgres password=milenov3790 dbname=test port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func testDbCreate(c *gin.Context) {
	dsn := "host=localhost user=postgres password=milenov3790 dbname=test port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func getProductDetails(c *gin.Context) {
	dsn := "host=localhost user=postgres password=milenov3790 dbname=test port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	id := c.Params.ByName("id")
	var product Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, product)
	}
}

func getProducts(c *gin.Context) {
	dsn := "host=localhost user=postgres password=milenov3790 dbname=test port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var product []Product
	if err := db.Find(&product).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, product)
	}
}
