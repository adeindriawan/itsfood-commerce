package controllers

import (
	// "fmt"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

type NewOrder struct {
	OrderedBy uint64 			`json:"ordered_by"`
	OrderedFor string			`json:"ordered_for"`
	OrderedTo string			`json:"ordered_to"`
	NumOfMenus int				`json:"num_of_menus"`
	QtyOfMenus int				`json:"qty_of_menus"`
	Amount int						`json:"amount"`
	Purpose string				`json:"purpose"`
	Activity string				`json:"activity"`
	SourceOfFund string		`json:"source_of_fund"`
	PaymentOption string	`json:"payment_option"`
	Info string						`json:"info"`
}

func _menuPreOrderHoursValidated(cartContent []Cart, orderedFor time.Time) bool {
	var minDeliveryTime time.Time
	var minHours, minDays uint
	minHours = 0
	minDays = 0
	now := time.Now()
	for _, search := range cartContent {
		if search.PreOrderHours > minHours {
			minHours = search.PreOrderHours
		}

		if search.PreOrderDays > minDays {
			minDays = search.PreOrderDays
		}

		if minDays > 0 {
			minDeliveryTime = now.AddDate(0, 0, int(minDays))
		} else {
			minDeliveryTime = now.Add(time.Hour * time.Duration(minHours))
		}
	}

	return !minDeliveryTime.After(orderedFor)
}

func CreateOrder(c *gin.Context) {
	// check apakah user sudah terautentikasi dan merupakan customer -> sudah teratasi dengan middleware
	// check apakan user/customer tersebut berstatus aktif -> sudah teratasi dengan middleware

	var order NewOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	u := c.MustGet("user").(models.User)
	userId := strconv.Itoa(int(u.ID))
	cartContent, errCartContent := GetUserCartContent(userId)
	if errCartContent != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errCartContent.Error(),
			"result": nil,
			"description": "Gagal mengambil isi keranjang user.",
		})
		return
	}
	// check apakah cart customer tersebut berisi setidaknya 1 item
	totalItems := GetUserCartItems(cartContent)
	if totalItems == 0 {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": nil,
			"result": nil,
			"description": "Gagal membuat order baru. Keranjang user masih kosong.",
		})
		return
	}

	order.NumOfMenus = totalItems
	order.QtyOfMenus = GetUserCartQty(cartContent)
	order.Amount = GetUserCartAmount(cartContent)
	
	// check apakah item di cartnya sudah memenuhi aturan min/max (sudah teratasi di controllers/cart) dan pre order
	orderedFor, errConvertingOrderedFor := time.Parse(time.RFC3339, order.OrderedFor)
	if errConvertingOrderedFor != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errConvertingOrderedFor.Error(),
			"result": nil,
			"description": "Gagal mengonversi waktu pengantaran pesanan.",
		})
		return
	}
	if !_menuPreOrderHoursValidated(cartContent, orderedFor) {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Order time ahead of the minimum delivery time.",
			"result": nil,
			"description": "Waktu pengantaran minimum lebih lama dibandingkan dengan waktu pengantaran pesanan.",
		})
		return
	}

	// Tambahkan record ke tabel orders dan order details
	// apakah ada order yang mengandung ITSMINE, jika ya, tembakkan ke API ITSMine
	// apakah ada menu yang vendornya memiliki default delivery cost/service charge, jika ya, tambahkan record ke costs
	// apakah ada menu yang vendornya memiliki telegram id, jika ya, kirim notifikasi telegram ke ID tersebut

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