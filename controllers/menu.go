package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
	"github.com/adeindriawan/itsfood-commerce/services"
)

type MenuData struct {
	ID uint										`json:"id"`
  Name string								`json:"name"`
	Description string				`json:"description"`
	VendorID uint							`json:"vendor_id"`
  VendorName string					`json:"vendor_name"`
	Type models.MenuCategory 	`json:"category"`
	RetailPrice uint 					`json:"retail_price"`
	WholesalePrice uint				`json:"wholesale_price"`
	COGS uint									`json:"cogs"`
	PreOrderDays uint 				`json:"pre_order_days"`
	PreOrderHours uint 				`json:"pre_order_hours"`
	MinOrderQty uint 					`json:"min_order_qty"`
	MaxOrderQty uint 					`json:"max_order_qty"`
	Image string 							`json:"image"`
	VendorNoteForMenus string	`json:"vendor_note_for_menus"`
	VendorDeliveryCost uint		`json:"vendor_delivery_cost"`
	VendorServiceCharge uint	`json:"vendor_service_charge"`
	VendorMinOrderAmount uint	`json:"vendor_min_order_amount"`
}

type MenuDataResponse struct {
	Data []MenuData		`json:"data"`
	RowsCount int64		`json:"rowsCount"`
	TotalRows int64		`json:"totalRows"`
	
}

func GetMenus(c *gin.Context) {
	var menu []MenuData
	var messages = []string{}
	params := c.Request.URL.Query()
	var categoryParam, doesCategoryParamExist = params["filters[category]"]
	var minOrderQtyParam, doesMinOrderQtyParamExist = params["filters[minOrderQty]"]
	var maxOrderQtyParam, doesMaxOrderQtyParamExist = params["filters[maxOrderQty]"]
	var preOrderDaysParam, doesPreOrderDaysParamExist = params["filters[preOrderDays]"]
	var preOrderHoursParam, doesPreOrderHoursParamExist = params["filters[preOrderHours]"]
	var priceMinParam, doesPriceMinParamExist = params["filters[price][min]"]
	var priceMaxParam, doesPriceMaxParamExist = params["filters[price][max]"]
	var lengthParam, doesLengthParamExist = params["length"]
	var pageParam, doesPageParamExist = params["page"]
	var searchParam, doesSearchParamExist = params["search"]
	var orderByParam, doesOrderByParamExist = params["orderBy"]
	var sortParam, doesSortParamExist = params["sort"]
	var vendorIdParam, doesVendorIdParamExist = params["vendorId"]
	query := services.DB.Table("menus m").
		Select(`m.id AS ID, m.name AS Name, m.description AS Description, v.id AS VendorID, u.name AS VendorName,
		m.type AS Type, m.retail_price AS RetailPrice, m.wholesale_price AS WholesalePrice, m.cogs AS COGS, m.pre_order_days AS PreOrderDays,
		m.pre_order_hours AS PreOrderHours, m.min_order_qty AS MinOrderQty, m.max_order_qty AS MaxOrderQty, m.image AS Image`).
		Joins("JOIN vendors v ON v.id = m.vendor_id").
		Joins("JOIN users u ON u.id = v.user_id").
		Where("m.status = ?", "Activated")
	if doesCategoryParamExist {
		menuCategory := categoryParam[0]
		switch menuCategory {
		case "Food", "Beverage", "Snack", "Fruit", "Grocery", "Others":
			query = query.Where("m.type = ?", menuCategory)	
		default:
			messages = append(messages, "Parameter Category yang diberikan tidak sesuai dengan kategori menu yang ada.")
		}
	}
	if doesMinOrderQtyParamExist {
		menuMinOrderQty := minOrderQtyParam[0]
		query = query.Where("m.min_order_qty <= ?", menuMinOrderQty)
	}
	if doesMaxOrderQtyParamExist {
		menuMaxOrderQty := maxOrderQtyParam[0]
		query = query.Where("m.max_order_qty >= ?", menuMaxOrderQty)
	}
	if doesPreOrderHoursParamExist {
		menuPreOrderHours := preOrderHoursParam[0]
		query = query.Where("m.pre_order_hours = ?", menuPreOrderHours)
	}
	if doesPreOrderDaysParamExist {
		menuPreOrderDays := preOrderDaysParam[0]
		query = query.Where("m.pre_order_days = ?", menuPreOrderDays)
	}
	if doesPriceMinParamExist {
		priceMin, err := strconv.Atoi(priceMinParam[0])
		if err != nil {
			messages = append(messages, "Parameter harga minimum tidak dapat dikonversi ke integer.")
		} else {
			query = query.Where("m.retail_price >= ?", priceMin)
		}
	}
	if doesPriceMaxParamExist {
		priceMax, err := strconv.Atoi(priceMaxParam[0])
		if err != nil {
			messages = append(messages, "Parameter harga maksimum tidak dapat dikonversi ke integer.")
		} else {
			query = query.Where("m.retail_price <= ?", priceMax)
		}
	}
	if doesSearchParamExist {
		search := searchParam[0]
		query = query.Where("m.name LIKE ?", "%" + search + "%").Or("u.name LIKE ?", "%" + search + "%")
	}
	if doesVendorIdParamExist {
		vendorId, err := strconv.Atoi(vendorIdParam[0])
		if err != nil {
			messages = append(messages, "Parameter ID vendor tidak dapat dikonversi ke integer.")
		} else {
			query = query.Where("v.id = ?", vendorId)
		}
	}
	if doesOrderByParamExist {
		var sort string
		if doesSortParamExist {
			sort = sortParam[0]
		} else {
			sort = "asc"
		}
		orderBy := orderByParam[0]
		if orderBy == "random" {
			query = query.Order("rand()")
		} else {
			query = query.Order(orderBy + " " + sort)
		}
	}
	query.Scan(&menu)
	totalRows := query.RowsAffected
	if doesLengthParamExist {
		length, err := strconv.Atoi(lengthParam[0])
		if err != nil {
			messages = append(messages, "Parameter Length tidak dapat dikonversi ke integer")
		} else {
			query = query.Limit(length)
		}
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
		fmt.Println(queryErr)
		messages = append(messages, "Ada kesalahan pada query")
		c.JSON(404, gin.H{
			"status": "failed",
			"errors": messages,
			"result": nil,
			"description": "Gagal mengambil data menu",
		})
	} else {
		menuData := &MenuDataResponse{
			Data: menu,
			RowsCount: rowsCount,
			TotalRows: totalRows,
		}

		c.JSON(200, gin.H{
			"status": "success",
			"errors": messages,
			"result": menuData,
			"description": "Berhasil mengambil data menu",
		})
	}
}

type MenuDetailsUri struct {
	Id int `uri:"id" binding:"required"`
}

func GetMenuDetails(c *gin.Context) {
	var messages = []string{}
	var menu MenuData
	var uri MenuDetailsUri
	if err := c.ShouldBindUri(&uri); err != nil {
		messages = append(messages, "Menu Id yang terkirim tidak valid.")
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": messages,
			"result": nil,
			"description": "Ada kesalahan terhadap nilai menu Id",
		})
		return
	} else {
		menuId := c.Param("id")
		query := services.DB.Table("menus m").
			Select(`m.id AS ID, m.name AS Name, m.description AS Description, v.id AS VendorID, u.name AS VendorName,
			v.vendor_note_for_menus AS VendorNoteForMenus, v.vendor_delivery_cost AS VendorDeliveryCost,
			v.vendor_service_charge AS VendorServiceCharge, v.vendor_min_order_amount AS VendorMinOrderAmount,
			m.type AS Type, m.retail_price AS RetailPrice, m.wholesale_price AS WholesalePrice, m.cogs AS COGS, m.pre_order_days AS PreOrderDays,
			m.pre_order_hours AS PreOrderHours, m.min_order_qty AS MinOrderQty, m.max_order_qty AS MaxOrderQty, m.image AS Image`).
			Joins("JOIN vendors v ON v.id = m.vendor_id").
			Joins("JOIN users u ON u.id = v.user_id").
			Where("m.status = ?", "Activated").
			Where("m.id =?", menuId).
			Order("m.id asc").Limit(1).Find(&menu)
		queryErr := query.Error
		queryRows := query.RowsAffected
		if queryErr != nil {
			messages = append(messages, "Ada kesalahan pada query.")
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": messages,
				"result": nil,
				"description": "Gagal mengambil data detail menu.",
			})
			return
		} else if queryRows == 0 {
			messages = append(messages, "Tidak ada menu aktif dengan ID tersebut.")
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": messages,
				"result": nil,
				"description": "Gagal mengambil data detail menu.",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"status": "success",
				"errors": messages,
				"result": menu,
				"description": "Berhasil mengambil data detail menu.",
			})
		}
	}
}
