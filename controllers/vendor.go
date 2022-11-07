package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/services"
)

type VendorInfoUri struct {
	Id uint	`uri:"id" binding:"required"`
}

type VendorInfo struct {
	Name string			`json:"name"`
	Address string	`json:"address"`
	NumOfMenus uint	`json:"num_of_menus"`
}

func GetVendorInfo(c *gin.Context) {
	var messages = []string{}
	var info VendorInfo
	var uri VendorInfoUri

	if err := c.ShouldBindUri(&uri); err != nil {
		messages = append(messages, err.Error())
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": messages,
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}
	vendorId := c.Param("id")
	subQuery := `
		SELECT (count(id)) FROM menus WHERE vendor_id = '` + vendorId + `' AND status = 'Activated'`
	query := services.DB.Table("vendors v").
	Select("u.name AS Name, v.address AS Address, ("+ subQuery +") AS NumOfMenus").
	Joins("JOIN users AS u ON v.user_id = u.id").
	Where("v.id = ?", vendorId).Order("v.id asc").Limit(1).Find(&info)
	queryError := query.Error
	if queryError != nil {
		messages = append(messages, queryError.Error())
		c.JSON(400, gin.H{
			"status": "failed",
			"messages": messages,
			"result": nil,
			"description": "Gagal mengambil data vendor.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": messages,
		"result": info,
		"description": "Berhasil mengambil info vendor.",
	})
}