package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
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
	_, errorSendingTelegram := services.SendTelegramToGroup("Chat test from Gin")
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