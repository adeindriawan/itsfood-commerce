package controllers

import (
	"runtime"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
	"github.com/adeindriawan/itsfood-commerce/utils"
	"strings"
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

func CreateOrderV1(c *gin.Context) {
	var order OrderPayload
	var errors = []string{}
	const itsmineVendorId int = 112
	var itsmineData []map[string]int
	itsmineOrder := make(map[string]int)

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(422, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	customerContext := c.MustGet("customer").(models.Customer)
	customerId := strconv.Itoa(int(customerContext.ID))
	cartContent, errCartContent := GetCustomerCartContent(customerId)
	if errCartContent != nil {
		c.JSON(500, gin.H{
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
		c.JSON(404, gin.H{
			"status": "failed",
			"errors": nil,
			"result": nil,
			"description": "Gagal membuat order baru. Keranjang user masih kosong.",
		})
		return
	}
	
	orderedFor, errConvertingOrderedFor := time.Parse(time.RFC3339, order.OrderedFor)
	if errConvertingOrderedFor != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"errors": errConvertingOrderedFor.Error(),
			"result": nil,
			"description": "Gagal mengonversi waktu pengantaran pesanan.",
		})
		return
	}
	isPreOrderValidated, minDeliveryTime := _menuPreOrderValidated(cartContent, orderedFor)
	if !isPreOrderValidated {
		c.JSON(422, gin.H{
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
		CreatedBy: customerContext.User.Name,
	}

	// Tambahkan record ke tabel orders dan order details
	creatingOrder := services.DB.Create(&newOrder)
	errorCreatingOrder := creatingOrder.Error
	if errorCreatingOrder != nil {
		c.JSON(512, gin.H{
			"status": "failed",
			"errors": errorCreatingOrder.Error(),
			"result": order,
			"description": "Gagal membuat menyimpan order baru ke dalam database.",
		})
		return
	}
	newOrderID := newOrder.ID
	year, month, day := newOrder.OrderedFor.Date()
	orderedForYear := strconv.Itoa(year)
	orderedForMonth := strconv.Itoa(int(month))
	orderedForDay := strconv.Itoa(day)
	orderedForDate := orderedForYear + "-" + orderedForMonth + "-" + orderedForDay

	for _, v := range cartContent {
		newOrderDetail := models.OrderDetail{
			OrderID: newOrderID,
			MenuID: v.MenuID,
			Qty: v.Qty,
			Price: v.Price,
			COGS: v.COGS,
			Note: v.Note,
			Status: "Ordered",
			CreatedAt: time.Now(),
			CreatedBy: customerContext.User.Name,
		}
		creatingOrderDetails := services.DB.Create(&newOrderDetail)
		errorCreatingOrderDetails := creatingOrderDetails.Error
		if errorCreatingOrderDetails != nil {
			c.JSON(512, gin.H{
				"status": "failed",
				"errors": errorCreatingOrderDetails.Error(),
				"result": v,
				"description": "Gagal menyimpan data detail order.",
			})
			return
		}

		if int(v.VendorID) == itsmineVendorId {
			itsmineOrder["id"] = int(newOrderDetail.ID)
			itsmineOrder["qty"] = int(v.Qty)
			itsmineOrder["product_id"] = int(v.MenuID)

			itsmineData = append(itsmineData, itsmineOrder)
		}
	}

	addVendorDefaultCosts(newOrderID)
	notifyVendorsViaTelegram(newOrderID)
	notifyTelegramGroup(newOrder, customerContext)
	notifyAdminsViaEmail(newOrder, customerContext, cartContent)

	orderParam := map[string]interface{}{
		"id": newOrderID,
		"ordered_to": newOrder.OrderedTo,
		"ordered_for": orderedForDate,
		"info": newOrder.Info,
	}
	customerParam := map[string]interface{}{
		"id": customerContext.ID,
		"name": customerContext.User.Name,
		"phone": customerContext.User.Phone,
		"unit_id": customerContext.Unit.ID,
		"unit_name": customerContext.Unit.Name,
	}
	params := map[string]interface{}{
		"order": orderParam,
		"customer": customerParam,
		"items": itsmineData,
	}

	if len(itsmineData) > 0 {
		_, errorSendingItsmineOrder := utils.AddItsmineOrder(params)
		if errorSendingItsmineOrder != nil {
			errors = append(errors, errorSendingItsmineOrder.Error())
		}
	}
	DestroyCustomerCart(customerId)

	c.JSON(201, gin.H{
		"status": "success",
		"result": map[string]interface{}{
			"order": newOrder,
		},
		"errors": errors,
		"description": "Berhasil membuat order baru.",
	})
}

func CreateOrder(c *gin.Context) {
	runtime.GOMAXPROCS(4)

	var order OrderPayload
	var errors = []string{}
	const itsmineVendorId int = 112
	var itsmineData []map[string]int
	itsmineOrder := make(map[string]int)

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(422, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	customerContext := c.MustGet("customer").(models.Customer)
	customerId := strconv.Itoa(int(customerContext.ID))
	cartContent, errCartContent := GetCustomerCartContent(customerId)
	if errCartContent != nil {
		c.JSON(500, gin.H{
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
		c.JSON(404, gin.H{
			"status": "failed",
			"errors": nil,
			"result": nil,
			"description": "Gagal membuat order baru. Keranjang user masih kosong.",
		})
		return
	}
	
	orderedFor, errConvertingOrderedFor := time.Parse(time.RFC3339, order.OrderedFor)
	if errConvertingOrderedFor != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"errors": errConvertingOrderedFor.Error(),
			"result": nil,
			"description": "Gagal mengonversi waktu pengantaran pesanan.",
		})
		return
	}
	isPreOrderValidated, minDeliveryTime := _menuPreOrderValidated(cartContent, orderedFor)
	if !isPreOrderValidated {
		c.JSON(422, gin.H{
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
		CreatedBy: customerContext.User.Name,
	}

	// Tambahkan record ke tabel orders dan order details
	creatingOrder := services.DB.Create(&newOrder)
	errorCreatingOrder := creatingOrder.Error
	if errorCreatingOrder != nil {
		c.JSON(512, gin.H{
			"status": "failed",
			"errors": errorCreatingOrder.Error(),
			"result": order,
			"description": "Gagal membuat menyimpan order baru ke dalam database.",
		})
		return
	}
	newOrderID := newOrder.ID
	year, month, day := newOrder.OrderedFor.Date()
	orderedForYear := strconv.Itoa(year)
	orderedForMonth := strconv.Itoa(int(month))
	orderedForDay := strconv.Itoa(day)
	orderedForDate := orderedForYear + "-" + orderedForMonth + "-" + orderedForDay

	for _, v := range cartContent {
		newOrderDetail := models.OrderDetail{
			OrderID: newOrderID,
			MenuID: v.MenuID,
			Qty: v.Qty,
			Price: v.Price,
			COGS: v.COGS,
			Note: v.Note,
			Status: "Ordered",
			CreatedAt: time.Now(),
			CreatedBy: customerContext.User.Name,
		}
		creatingOrderDetails := services.DB.Create(&newOrderDetail)
		errorCreatingOrderDetails := creatingOrderDetails.Error
		if errorCreatingOrderDetails != nil {
			c.JSON(512, gin.H{
				"status": "failed",
				"errors": errorCreatingOrderDetails.Error(),
				"result": v,
				"description": "Gagal menyimpan data detail order.",
			})
			return
		}

		if int(v.VendorID) == itsmineVendorId {
			itsmineOrder["id"] = int(newOrderDetail.ID)
			itsmineOrder["qty"] = int(v.Qty)
			itsmineOrder["product_id"] = int(v.MenuID)

			itsmineData = append(itsmineData, itsmineOrder)
		}
	}

	orderParam := map[string]interface{}{
		"id": newOrderID,
		"ordered_to": newOrder.OrderedTo,
		"ordered_for": orderedForDate,
		"info": newOrder.Info,
	}
	customerParam := map[string]interface{}{
		"id": customerContext.ID,
		"name": customerContext.User.Name,
		"phone": customerContext.User.Phone,
		"unit_id": customerContext.Unit.ID,
		"unit_name": customerContext.Unit.Name,
	}
	params := map[string]interface{}{
		"order": orderParam,
		"customer": customerParam,
		"items": itsmineData,
	}

	if len(itsmineData) > 0 {
		go utils.AddItsmineOrder(params)
	}

	go addVendorDefaultCosts(newOrderID)
	go notifyVendorsViaTelegram(newOrderID)
	go notifyTelegramGroup(newOrder, customerContext)
	go notifyAdminsViaEmail(newOrder, customerContext, cartContent)

	DestroyCustomerCart(customerId)

	c.JSON(201, gin.H{
		"status": "success",
		"result": map[string]interface{}{
			"order": newOrder,
		},
		"errors": errors,
		"description": "Berhasil membuat order baru.",
	})
}

func notifyVendorsViaTelegram(orderID uint64) {
	var order models.Order
	services.DB.Preload("OrderDetail").Where("id = ?", orderID).First(&order)
	orderDetails := order.OrderDetail
	var wasSent = []string{}
	for _, v := range orderDetails {
		orderDetailID := v.ID
		isSent := sendTelegramNotificationToVendor(orderDetailID)
		if isSent {
			wasSent = append(wasSent, "true")
		}
	}

	if len(wasSent) > 0 {
		var status string
		if len(wasSent) == len(orderDetails) {
			status = "ForwardedEntirely"
		} else {
			status = "ForwardedPartially"
		}
		models.UpdateOrder(map[string]interface{}{"id": orderID}, map[string]interface{}{"status": status, "updated_at": time.Now(), "created_by": "Itsfood Commerce System"})
	}
}

func sendTelegramNotificationToVendor(orderDetailID uint64) bool {
	var orderDetail models.OrderDetail
	services.DB.Preload("Order.Customer.Unit").Preload("Order.Customer.User").Preload("Menu.Vendor.User").Where("id = ?", orderDetailID).First(&orderDetail)
	vendorTelegramID := orderDetail.Menu.Vendor.VendorTelegramID
	if vendorTelegramID != "" && vendorTelegramID != "NULL" {
		orderID := strconv.Itoa(int(orderDetail.Order.ID))
		vendorName := orderDetail.Menu.Vendor.User.Name
		menuName := orderDetail.Menu.Name
		menuQty := strconv.Itoa(int(orderDetail.Qty))
		customerName := orderDetail.Order.Customer.User.Name
		unitName := orderDetail.Order.Customer.Unit.Name
		orderedForDate := utils.ConvertDateToPhrase(orderDetail.Order.OrderedFor, true)

		telegramMessage := "Ada order baru untuk " + vendorName + " dengan ID #" + orderID + " dari " + customerName + " dari " + unitName + " berupa " + menuName + " sebanyak " + menuQty + " porsi"
		telegramMessage += " untuk diantar pada " + orderedForDate

		models.UpdateOrderDetail(map[string]interface{}{"id": orderDetailID}, map[string]interface{}{"status": "Sent", "updated_at": time.Now(), "created_by": "Itsfood Commerce System"})

		_sendTelegramToVendor(telegramMessage, vendorTelegramID)
		return true
	}

	return false
}

func _sendTelegramToVendor(message string, chatID string) {
	services.SendTelegramToVendor(message, chatID)
}

func addVendorDefaultCosts(orderID uint64) {
	var order models.Order
	services.DB.Preload("OrderDetail").Where("id = ?", orderID).First(&order)
	orderDetails := order.OrderDetail
	for _, v := range orderDetails {
		orderDetailID := v.ID
		saveVendorDefaultCosts(orderDetailID)
	}
}

func saveVendorDefaultCosts(orderDetailID uint64) bool {
	var orderDetail models.OrderDetail
	services.DB.Preload("Menu.Vendor").Where("id = ?", orderDetailID).First(&orderDetail)
	vendorDeliveryCost := orderDetail.Menu.Vendor.VendorDeliveryCost
	vendorServiceCharge := orderDetail.Menu.Vendor.VendorServiceCharge
	if vendorDeliveryCost != 0 || vendorServiceCharge != 0 {
		if vendorDeliveryCost != 0 {
			newCost := models.Cost{
				OrderDetailID: orderDetail.ID,
				Amount: vendorDeliveryCost,
				Reason: "Delivery cost",
				Status: "Unpaid",
				CreatedAt: time.Now(),
				CreatedBy: "Itsfood Commerce System",
			}

			services.DB.Create(&newCost)
		}

		if vendorServiceCharge != 0 {
			newCost := models.Cost{
				OrderDetailID: orderDetail.ID,
				Amount: vendorServiceCharge,
				Reason: "Service charge",
				Status: "Unpaid",
				CreatedAt: time.Now(),
				CreatedBy: "Itsfood Commerce System",
			}

			services.DB.Create(&newCost)
		}

		return true
	}

	return false
}

func notifyTelegramGroup(newOrder models.Order, customerContext models.Customer) {
	_sendTelegramToGroup(newOrder, customerContext)
}

func _sendTelegramToGroup(newOrder models.Order, customerContext models.Customer) (bool, error) {
	telegramMessage := "Ada order baru nomor #"
	orderId := strconv.Itoa(int(newOrder.ID))
	orderedForDate := utils.ConvertDateToPhrase(newOrder.OrderedFor, true)
	telegramMessage += orderId + " dari " + customerContext.User.Name + " di " + customerContext.Unit.Name
	telegramMessage += " untuk diantar pada " + orderedForDate
	telegramMessage += ", klik <a href='https://itsfood.id/publics/view-order/" + orderId + "'> di sini</a> untuk detail."

	_, errorSendingTelegram := services.SendTelegramToGroup(telegramMessage)
	if errorSendingTelegram != nil {
		return false, errorSendingTelegram
	}

	return true, nil
}

func notifyAdminsViaEmail(newOrder models.Order, customerContext models.Customer, cartContent []Cart) {
	_sendEmailToAdmins(newOrder, customerContext, cartContent)
}

func _sendEmailToAdmins(newOrder models.Order, customerContext models.Customer, cartContent []Cart) (bool, error) {
	var admins []models.Admin
	query := services.DB.Preload("User").Find(&admins)
	queryError := query.Error
	if queryError != nil {
		return false, queryError
	}

	for _, v := range admins {
		emailBody := _newOrderEmailBody(newOrder, customerContext, cartContent, v.ID)
		services.SendMail(v.Email, "[Itsfood] Pesanan Baru", emailBody)
	}

	return true, nil
}

func _newOrderEmailBody(newOrder models.Order, customerContext models.Customer, cartContent []Cart, adminID uint64) string {
	adminId := strconv.Itoa(int(adminID))
	orderId := strconv.Itoa(int(newOrder.ID))
	orderedForDate := utils.ConvertDateToPhrase(newOrder.OrderedFor, true)
	customerCartAmount := strconv.Itoa(int(GetCustomerCartAmount(cartContent)))
	cartDetails := _cartDetailsForEmail(cartContent)
	emailMessage := "Dear admin Itsfood,<br><br><br>"
	emailMessage += "Ada order baru dengan ID #" + orderId + " dari " + customerContext.User.Name + " di " + customerContext.Unit.Name
	emailMessage += " dengan rincian sebagai berikut:<br>"
	emailMessage += cartDetails + "<br>"
	emailMessage += "Total penjualan: Rp" + customerCartAmount + "<br>"
	emailMessage += "Diantar pada: " + orderedForDate + "<br>"
	emailMessage += "Tujuan: " + newOrder.OrderedTo + "<br>"
	emailMessage += "Untuk keperluan: " + newOrder.Purpose + "<br>"
	emailMessage += "Informasi tambahan: " + newOrder.Info + "<br>"
	emailMessage += "Kontak pembeli: " + customerContext.User.Phone + "<br>"
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

type OrderResult struct {
	ID uint64							`json:"id"`
	OrderedFor string			`json:"ordered_for"`
	OrderedTo string 			`json:"ordered_to"`
	NumOfMenus uint				`json:"num_of_menus"`
	QtyOfMenus uint 			`json:"qty_of_menus"`
	Amount uint64					`json:"amount"`
	Purpose string				`json:"purpose"`
	Activity string 			`json:"activity"`
	SourceOfFund string		`json:"source_of_fund"`
	PaymentOption string	`json:"payment_option"`
	Info string						`json:"info"`
	Status string					`json:"status"`
	CreatedAt string			`json:"created_at"`
}

func GetOrders(c *gin.Context) {
	var orders []OrderResult
	var messages = []string{}

	customerContext := c.MustGet("customer").(models.Customer)
	customerID := customerContext.ID

	params := c.Request.URL.Query()
	idParam, doesIdParamExist := params["id"]
	paidParam, doesPaidParamExist := params["paid"]
	lengthParam, doesLengthParamExist := params["length"]
	pageParam, doesPageParamExist := params["page"]
	statusParam, doesStatusParamExist := params["status"]
	startOrderDateParam, doesStartOrderDateParamExist := params["order_date[start]"]
	endOrderDateParam, doesEndOrderDateParamExist := params["order_date[end]"]
	startDeliveryDateParam, doesStartDeliveryDateParamExist := params["delivery_date[start]"]
	endDeliveryDateParam, doesEndDeliveryDateParamExist := params["delivery_date[end]"]
	purposeParam, doesPurposeParamExist := params["purpose"]
	orderQuery := services.DB.Table("orders").
		Select(`
			id AS ID, ordered_for AS OrderedFor, ordered_to AS OrderedTo, num_of_menus AS NumOfMenus,
			qty_of_menus AS QtyOfMenus, amount AS Amount, purpose AS Purpose, activity AS Activity,
			source_of_fund AS SourceOfFund, payment_option AS PaymentOption, info AS Info,
			status AS Status, created_at AS CreatedAt
		`).
		Where("ordered_by = ?", customerID)
	
	if doesStatusParamExist {
		status := statusParam[0]
		orderQuery = orderQuery.Where("status IN ?", strings.Split(status, ","))
	}

	if doesStartOrderDateParamExist && !doesEndOrderDateParamExist {
		startOrderDate := startOrderDateParam[0]
		orderQuery = orderQuery.Where("created_at >= ?", startOrderDate)
	}

	if !doesStartOrderDateParamExist && doesEndOrderDateParamExist {
		endOrderDate := endOrderDateParam[0]
		orderQuery = orderQuery.Where("created_at <= ?", endOrderDate)
	}

	if doesStartOrderDateParamExist && doesEndOrderDateParamExist {
		startOrderDate := startOrderDateParam[0]
		endOrderDate := endOrderDateParam[0]
		orderQuery = orderQuery.Where("created_at BETWEEN ? AND ?", startOrderDate, endOrderDate)
	}

	if doesStartDeliveryDateParamExist && !doesEndDeliveryDateParamExist {
		startDeliveryDate := startDeliveryDateParam[0]
		orderQuery = orderQuery.Where("ordered_for >= ?", startDeliveryDate)
	}

	if !doesStartDeliveryDateParamExist && doesEndDeliveryDateParamExist {
		endDeliveryDate := endDeliveryDateParam[0]
		orderQuery = orderQuery.Where("ordered_for <= ?", endDeliveryDate)
	}

	if doesStartDeliveryDateParamExist && doesEndDeliveryDateParamExist {
		startDeliveryDate := startDeliveryDateParam[0]
		endDeliveryDate := endDeliveryDateParam[0]
		orderQuery = orderQuery.Where("ordered_for BETWEEN ? AND ?", startDeliveryDate, endDeliveryDate)
	}

	if doesPurposeParamExist {
		purpose := purposeParam[0]
		orderQuery = orderQuery.Where("purpose LIKE ?", "%"+purpose+"%")
	}
	
	if doesIdParamExist {
		id := idParam[0]
		orderQuery = orderQuery.Where("id = ?", id)
	}
	if doesPaidParamExist {
		paid := paidParam[0]
		if paid == "false" {
			orderQuery = orderQuery.Where("status IN ?", []string{"Created", "ForwardedPartially", "ForwardedEntirely", "Processed", "Completed", "BilledPartially", "BilledEntirely"})
		}

		if paid == "true" {
			orderQuery = orderQuery.Where("status IN ?", []string{"Paid", "PaidAndBilledPartially", "PaidAndBilledEntirely", "PaidByCustomerAndToVendor"})
		}
	}
	orderQuery.Scan(&orders)
	totalRows := orderQuery.RowsAffected
	if doesLengthParamExist {
		length, err := strconv.Atoi(lengthParam[0])
		if err != nil {
			messages = append(messages, "Parameter Length tidak dapat dikonversi ke integer")
		} else {
			orderQuery = orderQuery.Limit(length)
		}
	}
	if doesPageParamExist {
		if doesLengthParamExist {
			page, _ := strconv.Atoi(pageParam[0])
			length, _ := strconv.Atoi(lengthParam[0])
			offset := (page - 1) * length
			orderQuery = orderQuery.Offset(offset)
		} else {
			messages = append(messages, "Tidak ada parameter Length, maka parameter Page diabaikan.")
		}
	}
	orderQuery.Scan(&orders)
	rowsCount := orderQuery.RowsAffected

	if orderQuery.Error != nil {
		c.JSON(512, gin.H{
			"status": "failed",
			"errors": orderQuery.Error.Error(),
			"result": nil,
			"description": "Gagal mengeksekusi query.",
		})
	}

	orderData := map[string]interface{}{
		"data": orders,
		"rows_count": rowsCount,
		"total_rows": totalRows,
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": orderData,
		"errors": messages,
		"description": "Berhasil mengambil data order dari user ini.",
	})
}

type OrderDetailsUri struct {
	Id uint64 `uri:"id" binding:"required"`
}

type OrderDetailResult struct {
	ID uint64					`json:"id"`
	MenuID uint64			`json:"menu_id"`
	MenuName string		`json:"menu_name"`
	MenuImage string 	`json:"menu_image"`
	VendorName string	`json:"vendor_name"`
	Price uint64			`json:"price"`
	Qty uint					`json:"qty"`
	Note string				`json:"note"`
	Status string			`json:"status"`
}

func OrderDetails(c *gin.Context) {
	var order OrderResult
	var orderDetail []OrderDetailResult
	var uri OrderDetailsUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
		})
		return
	}

	orderId := c.Param("id")
	orderQuery := services.DB.Table("orders").
				Select(`id AS ID, ordered_for AS OrderedFor, ordered_to AS OrderedTo, 
				num_of_menus AS NumOfMenus, qty_of_menus AS QtyOfMenus, amount AS Amount, purpose AS Purpose, 
				activity AS Activity, 
				source_of_fund AS SourceOfFund, payment_option AS PaymentOption, info AS Info, 
				status AS Status, created_at AS CreatedAt`).
				Where("id = ?", orderId).Order("id ASC").Limit(1).Scan(&order)
	orderQueryError := orderQuery.Error

	orderDetailQuery := services.DB.Table("order_details od").
				Select(`od.id AS ID, od.menu_id AS MenuID, m.name AS MenuName, m.image AS MenuImage, 
				u.name AS VendorName, od.price AS Price, od.qty AS Qty,
				od.note AS Note, od.status AS Status`).
				Joins("JOIN menus m ON od.menu_id = m.id").
				Joins("JOIN vendors v ON m.vendor_id = v.id").
				Joins("JOIN users u ON v.user_id = u.id").
				Where("order_id", orderId).Scan(&orderDetail)
	orderDetailQueryError := orderDetailQuery.Error
	
	if orderQueryError != nil || orderDetailQueryError != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"result": nil,
			"errors": orderQueryError.Error() + " & " + orderDetailQueryError.Error(),
			"description": "Gagal mengambil data order serta rinciannya dari database.",
		})
		return
	}

	result := map[string]interface{}{
		"order": order,
		"details": orderDetail,
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": result,
		"description": "Berhasil mengambil data order serta rinciannya.",
	})
}