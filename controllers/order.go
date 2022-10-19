package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func CreateOrder(c *gin.Context) {
	// check apakah user sudah terautentikasi dan merupakan customer
	// check apakan user/customer tersebut berstatus aktif
	// check apakah cart customer tersebut berisi setidaknya 1 item
	// check apakah item di cartnya sudah memenuhi aturan min/max dan pre order
	// Tambahkan record ke tabel orders dan order details
	// apakah ada order yang mengandung ITSMINE, jika ya, tembakkan ke API ITSMine
	// apakah ada menu yang vendornya memiliki default delivery cost/service charge, jika ya, tambahkan record ke costs
	// apakah ada menu yang vendornya memiliki telegram id, jika ya, kirim notifikasi telegram ke ID tersebut
	

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(200, gin.H{
			"status": "failed",
			"errors": err.Error(),
		})
		return
	}

	_, errorSendingTelegram := services.SendTelegramToGroup("Chat test from Gin")
	if errorSendingTelegram != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorSendingTelegram.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"result": order,
	})
}

type OrderDetailsUri struct {
	Id uint64 `uri:"id" binding:"required"`
}

func OrderDetails(c *gin.Context) {
	var order models.Order
	var uri OrderDetailsUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
		})
		return
	}

	orderId := c.Param("id")
	query := services.DB.Preload("Customer.User").Table("orders").Where("id = ?", orderId).Order("id ASC").Limit(1).Find(&order)
	queryError := query.Error
	
	if queryError != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": queryError.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": order,
	})
}