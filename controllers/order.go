package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
)

func CreateOrder(c *gin.Context) {
	// check apakah user sudah terautentikasi dan merupakan customer
	// check apakah cart customer tersebut berisi setidaknya 1 item
	// check apakah item di cartnya sudah memenuhi aturan min/max dan pre order
	// apakah ada order yang mengandung ITSMINE
	// apakah ada menu yang vendornya memiliki default delivery cost/service charge
	// apakah ada menu yang vendornya memiliki telegram id
	_, errorSendingTelegram := services.SendTelegram("Chat test from Gin")
	if errorSendingTelegram != nil {
		c.JSON(400, gin.H{
			"status": "failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}