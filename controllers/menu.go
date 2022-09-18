package controllers

import (
	"net/http"
	"fmt"
	"strconv"

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
	ID uint										`json:"id"`
  Name string								`json:"name"`
	Description string				`json:"description"`
	VendorID uint							`json:"vendor_id"`
  VendorName string					`json:"vendor_name"`
	Type models.MenuCategory 	`json:"type"`
	RetailPrice uint 					`json:"retail_price"`
	WholesalePrice uint				`json:"wholesale_price"`
	PreOrderDays uint 				`json:"pre_order_days"`
	PreOrderHours uint 				`json:"pre_order_hours"`
	MinOrderQty uint 					`json:"min_order_qty"`
	MaxOrderQty uint 					`json:"max_order_qty"`
	Image string 							`json:"image"`
}

type MenuDataResponse struct {
	Data []MenuData			`json:"data"`
	RowsCount int64			`json:"rowsCount"`
	RowsFiltered int64	`json:"rowsFiltered"`
}

func GetMenus(c *gin.Context) {
	var menu []MenuData
	var messages = []string{}
	params := c.Request.URL.Query()
	var typeParam, doesTypeParamExist = params["type"]
	var lengthParam, doesLengthParamExist = params["length"]
	var pageParam, doesPageParamExist = params["page"]
	var searchParam, doesSearchParamExist = params["search"]
	var sortParam, doesSortParamExist = params["sort_by"]
	query := models.DB.Table("menus m").
		Select(`m.id AS ID, m.name AS Name, m.description AS Description, v.id AS VendorID, u.name AS VendorName,
		m.type AS Type, m.retail_price AS RetailPrice, m.wholesale_price AS WholesalePrice, m.pre_order_days AS PreOrderDays,
		m.pre_order_hours AS PreOrderHours, m.min_order_qty AS MinOrderQty, m.max_order_qty AS MaxOrderQty, m.image AS Image`).
		Joins("JOIN vendors v ON v.id = m.vendor_id").
		Joins("JOIN users u ON u.id = v.user_id").
		Where("m.status = ?", "Activated")
	if doesTypeParamExist {
		menuType := typeParam[0]
		switch menuType {
		case "Food", "Beverage", "Snack", "Fruit", "Grocery", "Others":
			query = query.Where("m.type = ?", menuType)	
		default:
			messages = append(messages, "Parameter Type yang diberikan tidak sesuai dengan kategori menu yang ada.")
		}
	}
	if doesSearchParamExist {
		search := searchParam[0]
		query = query.Where("m.name LIKE ?", "%" + search + "%")
	}
	if doesSortParamExist {
		sort := sortParam[0]
		query = query.Order(sort + " asc")
	}
	if doesLengthParamExist {
		length, _ := strconv.Atoi(lengthParam[0])
		query = query.Limit(length)
	}
	if doesPageParamExist {
		if doesLengthParamExist {
			page, _ := strconv.Atoi(pageParam[0])
			length, _ := strconv.Atoi(lengthParam[0])
			offset := (page - 1) * length
			query = query.Offset(offset)
		} else {
			messages = append(messages, "Tidak ada parameter Length, maka parameter Page diabaikan.")
		}
	}
	query.Scan(&menu)

	rowsCount := query.RowsAffected
	queryErr := query.Error

	if queryErr != nil {
		c.AbortWithStatus(404)
		fmt.Println(queryErr)
	} else {
		menuData := &MenuDataResponse{
			Data: menu,
			RowsCount: rowsCount,
			RowsFiltered: rowsCount,
		}

		c.JSON(200, gin.H{
			"status": "success",
			"messages": messages,
			"result": menuData,
			"description": "Data menu berhasil diambil",
		})
		fmt.Println(rowsCount)
	}
}
