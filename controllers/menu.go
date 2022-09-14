package controllers

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
)

func SimpleRequest0(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "this is homepage",
	})
}

func SimpleRequest1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func SimpleRequest2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func SimpleRequest3(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func SimpleRequest4(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	message := name + " is " + action

	c.String(http.StatusOK, message)
}

func SimpleRequest5(c *gin.Context) {
	b := c.FullPath() == "/user/:name/*action" // true
	c.String(http.StatusOK, "%t", b)
}

type MenuData struct {
	ID uint							`json:"id"`
  Name string					`json:"name"`
	Description string	`json:"description"`
	VendorID uint				`json:"vendor_id"`
  VendorName string		`json:"vendor_name"`
	Type string 				`json:"type"`
	RetailPrice uint 		`json:"retail_price"`
	WholesalePrice uint	`json:"wholesale_price"`
	PreOrderDays uint 	`json:"pre_order_days"`
	PreOrderHours uint 	`json:"pre_order_hours"`
	MinOrderQty uint 		`json:"min_order_qty"`
	MaxOrderQty uint 		`json:"max_order_qty"`
	Image string 				`json:"image"`
}

func GetMenus(c *gin.Context) {
	var menu []MenuData
	if err := models.DB.Table("menus m").
							Select(
								`m.id AS ID, m.name AS Name, m.description AS Description, v.id AS VendorID, u.name AS VendorName,
								m.type AS Type, m.retail_price AS RetailPrice, m.wholesale_price AS WholesalePrice, m.pre_order_days AS PreOrderDays,
								m.pre_order_hours AS PreOrderHours, m.min_order_qty AS MinOrderQty, m.max_order_qty AS MaxOrderQty, m.image AS Image`).
							Joins("JOIN vendors v ON v.id = m.vendor_id").
							Joins("JOIN users u ON u.id = v.user_id").
							Limit(10).Scan(&menu).Error; err != nil {
								c.AbortWithStatus(404)
								fmt.Println(err)
							} else {
								c.JSON(200, menu)
							}
}
