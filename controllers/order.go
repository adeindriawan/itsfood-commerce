package controllers

import (
	"fmt"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

type OrderPayload struct {
	OrderedBy uint64 			`json:"ordered_by"`
	OrderedFor string			`json:"ordered_for"`
	OrderedTo string			`json:"ordered_to"`
	Purpose string				`json:"purpose"`
	Activity string				`json:"activity"`
	SourceOfFund string		`json:"source_of_fund"`
	PaymentOption string	`json:"payment_option"`
	Info string						`json:"info"`
}

func _menuPreOrderValidated(cartContent []Cart, orderedFor time.Time) (bool, time.Time) {
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
			minDate := now.AddDate(0, 0, int(minDays))
			minDeliveryTime = time.Date(minDate.Year(), minDate.Month(), minDate.Day(), 0, 0, 0, 0, minDate.Location())
		} else {
			minDeliveryTime = now.Add(time.Hour * time.Duration(minHours))
		}
	}

	return !minDeliveryTime.After(orderedFor), minDeliveryTime
}

func CreateOrder(c *gin.Context) {
	var order OrderPayload
	var errors = []string{}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	userContext := c.MustGet("user").(models.User)
	customerContext := c.MustGet("customer").(models.Customer)
	customerId := strconv.Itoa(int(customerContext.ID))
	cartContent, errCartContent := GetCustomerCartContent(customerId)
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
	totalItems := GetCustomerCartItems(cartContent)
	if totalItems == 0 {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": nil,
			"result": nil,
			"description": "Gagal membuat order baru. Keranjang user masih kosong.",
		})
		return
	}
	
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
	isPreOrderValidated, minDeliveryTime := _menuPreOrderValidated(cartContent, orderedFor)
	if !isPreOrderValidated {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Order time ahead of the minimum delivery time.",
			"result": map[string]interface{}{
				"Minimum delivery time": minDeliveryTime,
			},
			"description": "Waktu pengantaran pesanan yang diinginkan lebih cepat dibandingkan dengan waktu minimum pengiriman.",
		})
		return
	}

	orderedBy := customerContext.ID
	newOrder := models.Order{
		OrderedBy: orderedBy,
		OrderedFor: orderedFor,
		OrderedTo: order.OrderedTo,
		NumOfMenus: uint(totalItems),
		QtyOfMenus: uint(GetCustomerCartQty(cartContent)),
		Amount: uint64(GetCustomerCartAmount(cartContent)),
		Purpose: order.Purpose,
		Activity: order.Activity,
		SourceOfFund: order.SourceOfFund,
		PaymentOption: order.PaymentOption,
		Info: order.Info,
		Status: "Created",
		CreatedAt: time.Now(),
		CreatedBy: userContext.Name,
	}

	// Tambahkan record ke tabel orders dan order details
	creatingOrder := services.DB.Create(&newOrder)
	errorCreatingOrder := creatingOrder.Error
	if errorCreatingOrder != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingOrder.Error(),
			"result": order,
			"description": "Gagal membuat menyimpan order baru.",
		})
		return
	}
	newOrderID := newOrder.ID
	var orderModel models.Order
	getOrderData := services.DB.Preload("Customer.User").Preload("Customer.Unit").Find(&orderModel, newOrderID)
	errorGettingOrderData := getOrderData.Error
	if errorGettingOrderData != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorGettingOrderData.Error(),
			"result": nil,
			"description": "Gagal mengambil data order.",
		})
		return
	}

	for _, v := range cartContent {
		newOrderDetails := models.OrderDetail{
			OrderID: newOrderID,
			MenuID: v.MenuID,
			Qty: v.Qty,
			Price: v.Price,
			COGS: v.COGS,
			Status: "Ordered",
			CreatedAt: time.Now(),
			CreatedBy: userContext.Name,
		}
		creatingOrderDetails := services.DB.Create(&newOrderDetails)
		errorCreatingOrderDetails := creatingOrderDetails.Error
		if errorCreatingOrderDetails != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": errorCreatingOrderDetails.Error(),
				"result": v,
				"description": "Gagal menyimpan data detail order.",
			})
			return
		}
	}
	DestroyCustomerCart(customerId)

	// apakah ada order yang mengandung ITSMINE, jika ya, tembakkan ke API ITSMine
	// apakah ada menu yang vendornya memiliki default delivery cost/service charge, jika ya, tambahkan record ke costs
	// apakah ada menu yang vendornya memiliki telegram id, jika ya, kirim notifikasi telegram ke ID tersebut

	_, errorSendingTelegram := _sendTelegramToGroup(orderModel)
	if errorSendingTelegram != nil {
		errors = append(errors, "Gagal mengirim notifikasi telegram order baru ke grup.")
	}
	
	_, errorSendingNewOrderEmail := _sendEmailToAdmins(orderModel, cartContent)
	if errorSendingNewOrderEmail != nil {
		errors = append(errors, "Gagal mengirim notifikasi email order baru ke admin.")
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": map[string]interface{}{
			"order": orderModel,
		},
		"errors": errors,
		"description": "Berhasil membuat order baru.",
	})
}

func _sendTelegramToGroup(orderModel models.Order) (bool, error) {
	telegramMessage := "Ada order baru nomor #"
	orderId := strconv.Itoa(int(orderModel.ID))
	orderedForYear := strconv.Itoa(orderModel.OrderedFor.Year())
	orderedForMonth := orderModel.OrderedFor.Month().String()
	orderedForDay := strconv.Itoa(orderModel.OrderedFor.Day())
	orderedForHour := strconv.Itoa(orderModel.OrderedFor.Hour())
	orderedForMinute := strconv.Itoa(orderModel.OrderedFor.Minute())
	telegramMessage += orderId + " dari " + orderModel.Customer.User.Name + " di " + orderModel.Customer.Unit.Name
	telegramMessage += " untuk diantar pada " + orderedForDay + " " + orderedForMonth + " " + orderedForYear
	telegramMessage += " " + orderedForHour + ":" + orderedForMinute
	telegramMessage += ", klik <a href='https://itsfood.id/publics/view-order/" + orderId + "'> di sini</a> untuk detail."

	_, errorSendingTelegram := services.SendTelegramToGroup(telegramMessage)
	if errorSendingTelegram != nil {
		return false, errorSendingTelegram
	}

	return true, nil
}

func _sendEmailToAdmins(orderModel models.Order, cartContent []Cart) (bool, error) {
	var admins []models.Admin
	query := services.DB.Preload("User").Find(&admins)
	queryError := query.Error
	if queryError != nil {
		return false, queryError
	}

	for _, v := range admins {
		emailBody := _newOrderEmailBody(orderModel, cartContent, v.ID)
		services.SendMail(v.Email, "[Itsfood] Pesanan Baru", emailBody)
		fmt.Println(v.Email)
	}

	return true, nil
}

func _newOrderEmailBody(orderModel models.Order, cartContent []Cart, adminID uint64) string {
	adminId := strconv.Itoa(int(adminID))
	orderId := strconv.Itoa(int(orderModel.ID))
	orderedForYear := strconv.Itoa(orderModel.OrderedFor.Year())
	orderedForMonth := orderModel.OrderedFor.Month().String()
	orderedForDay := strconv.Itoa(orderModel.OrderedFor.Day())
	orderedForHour := strconv.Itoa(orderModel.OrderedFor.Hour())
	orderedForMinute := strconv.Itoa(orderModel.OrderedFor.Minute())
	customerCartAmount := strconv.Itoa(int(GetCustomerCartAmount(cartContent)))
	cartDetails := _cartDetailsForEmail(cartContent)
	emailMessage := "Dear admin Itsfood,<br><br><br>"
	emailMessage += "Ada order baru dengan ID #" + orderId + " dari " + orderModel.Customer.User.Name + " di " + orderModel.Customer.Unit.Name
	emailMessage += " dengan rincian sebagai berikut:<br>"
	emailMessage += cartDetails + "<br>"
	emailMessage += "Total penjualan: Rp" + customerCartAmount + "<br>"
	emailMessage += "Diantar pada: " + orderedForDay + " " + orderedForMonth + " " + orderedForYear + " " + orderedForHour + ":" + orderedForMinute + "<br>"
	emailMessage += "Tujuan: " + orderModel.OrderedTo + "<br>"
	emailMessage += "Untuk keperluan: " + orderModel.Purpose + "<br>"
	emailMessage += "Informasi tambahan: " + orderModel.Info + "<br>"
	emailMessage += "Kontak pembeli: " + orderModel.Customer.User.Phone + "<br>"
	emailMessage += "Silakan klik link di bawah ini untuk memproses:<br>"
	emailMessage += "<a style='font-size:14px; font-weight:bold; text-decoration:none; line-height:40px; width:100%; display:inline-block;' href='https://itsfood.id/publics/proceed-order/"+ orderId +"/" + adminId +"'><span style='color:#000091'>Proses Pesanan Ini</span></a>"

	return emailMessage
}

func _cartDetailsForEmail(customerCartContent []Cart) string {
	cartDetails := "<table><thead><tr><th>Nama Menu</th><th>Nama Vendor</th><th>Jumlah</th><th>Harga Pokok</th><th>Harga Jual</th><th>Pembayaran ke Vendor</th><th>Pembayaran dari Pembeli</th></tr></thead><tbody>"
	for _, v := range(customerCartContent) {
		Qty := strconv.Itoa(int(v.Qty))
		COGS := strconv.Itoa(int(v.COGS))
		Price := strconv.Itoa(int(v.Price))
		paymentToVendorNominal := int(v.Qty) * int(v.COGS)
		paymentToVendor := strconv.Itoa(int(paymentToVendorNominal))
		paymentFromCustomerNominal := int(v.Qty) * int(v.Price)
		paymentFromCustomer := strconv.Itoa(int(paymentFromCustomerNominal))
		cartDetails += "<tr><td>" + v.Name + "</td><td>" + v.VendorName + "</td><td>" + Qty + " porsi</td></td>Rp" + COGS + "</td><td>Rp" + Price + "</td><td>Rp" + paymentToVendor + "</td></td>Rp" + paymentFromCustomer + "</td></tr>"
	}
	cartDetails += "</tbody></table>"

	return cartDetails
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