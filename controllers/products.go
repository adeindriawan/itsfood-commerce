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
	var create CreateProductInput
	if err := c.ShouldBindJSON(&create); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{Code: create.Code, Price: create.Price}
	models.DB.Create(&product)
	fmt.Println(create)
	c.JSON(http.StatusOK, gin.H{"data": product})
}

type UpdateProductInput struct {
	Code	string	`json:"code"`
	Price	uint		`json:"price"`
}

func UpdateProduct(c *gin.Context) {
	//get model if exist
	var product models.Product
	if err := models.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak ditemukan!"})
		return
	}
	// validate input
	var update UpdateProductInput
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(product)
	fmt.Println(update)

	updatedProduct := models.Product{Code: update.Code, Price: update.Price}
	models.DB.Model(&product).Updates(&updatedProduct)
	c.JSON(http.StatusOK, gin.H{"data": product})
}