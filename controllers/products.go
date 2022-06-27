package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func FindProducts(c *gin.Context) {
	var products []models.Product
	if err := models.DB.Find(&products).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": products})
	}
}

func ProductDetails(c *gin.Context) {
	id := c.Params.ByName("id")
	var product models.Product
	if err := models.DB.Where("id = ?", id).First(&product).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": product})
	}
}

type CreateProductInput struct {
	Code 	string	`json:"code" binding:"required"`
	Price uint 		`json:"price" binding:"required"`
}

func CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{Code: input.Code, Price: input.Price}
	models.DB.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}