package controllers

import (
	// "fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/twinj/uuid"
)

type Cart struct {
	ID string 			`json:"id"`
	MenuID uint64 	`json:"menu_id"`
	Price uint64		`json:"price"`
	COGS uint64			`json:"cogs"`
	Subtotal uint64	`json:"subtotal"`
	VendorID uint64	`json:"vendor_id"`
	Qty uint64 			`json:"qty"`
}

func AuthCheck(c *gin.Context) (uint64, error) {
	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return 0, err
	}
	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func _cartContent(userId string) ([]Cart, error) {
	var content []Cart
	cart, err := services.GetRedis().Get("cart" + userId).Result()
	if err != nil {
		return nil, err
	}
	errUnmarshal := json.Unmarshal([]byte(cart), &content)
	if errUnmarshal != nil {
		return nil, err
	}
	return content, nil
}

func _menuExistsAndChangeQty(menuId uint64, qty uint64, cartContent []Cart) (bool, []Cart) {
	for i, search := range cartContent {
		if search.MenuID == menuId {
			cartContent[i].Qty += qty
			return true, cartContent
		}
	}
	return false, nil
}

func AddToCart(c *gin.Context) {
	var cart Cart
	var cartContent []Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(423, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Data yang masuk tidak dapat diproses.",
		})
		return
	}

	user, err := AuthCheck(c)
	if err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil token user yang ada. Unauthorized.",
		})
		return
	}
	userId := strconv.Itoa(int(user))
	currentCartContent, noCurrentCartContent := _cartContent(userId)
	if noCurrentCartContent != nil {
		cart.ID = uuid.NewV4().String()
		cartContent = append(cartContent, cart)
		json, err := json.Marshal(cartContent)
		if err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Gagal mengubah data ke dalam format JSON.",
			})
			return
		}
		errSave := services.GetRedis().Set("cart" + userId, json, 0).Err()
		if errSave != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": errSave.Error(),
				"result": nil,
				"description": "Tidak dapat menyimpan ke dalam keranjang belanja.",
			})
			return
		}
	} else {
		menuCheck, cartContent := _menuExistsAndChangeQty(cart.MenuID, cart.Qty, currentCartContent)
		if menuCheck {
			json, err := json.Marshal(cartContent)
			if err != nil {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": err.Error(),
					"result": nil,
					"description": "Gagal mengubah data ke dalam format JSON.",
				})
				return
			}
			errSave := services.GetRedis().Set("cart" + userId, json, 0).Err()
			if errSave != nil {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": errSave.Error(),
					"result": nil,
					"description": "Tidak dapat menyimpan ke dalam keranjang belanja.",
				})
				return
			}
		} else {
			cart.ID = uuid.NewV4().String()
			newCartContent := append(currentCartContent, cart)
			json, err := json.Marshal(newCartContent)
			if err != nil {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": err.Error(),
					"result": nil,
					"description": "Gagal mengubah data ke dalam format JSON.",
				})
				return
			}
			errSave := services.GetRedis().Set("cart" + userId, json, 0).Err()
			if errSave != nil {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": errSave.Error(),
					"result": nil,
					"description": "Tidak dapat menyimpan ke dalam keranjang belanja.",
				})
				return
			}
		}
	}
	
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": cart,
		"description": "Berhasil memasukkan menu ke dalam keranjang belanja.",
	})
}

func ViewCart(c *gin.Context) {
	user, err := AuthCheck(c)
	if err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil token user yang ada. Unauthorized.",
		})
		return
	}
	userId := strconv.Itoa(int(user))
	cart, err := _cartContent(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil data dalam keranjang belanja.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": cart,
		"description": "Berhasil mengambil data dari keranjang belanja.",
	})
}

func UpdateCart(c *gin.Context) {
	var cart Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Data yang masuk tidak dapat diproses.",
		})
		return
	}
	user, err := AuthCheck(c)
	if err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil token user yang ada. Unauthorized.",
		})
		return
	}
	userId := strconv.Itoa(int(user))
	cartContent, err := _cartContent(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak ada isi keranjang belanja dari user yang bersangkutan.",
		})
		return
	}
	for i, search := range cartContent {
		if search.ID == cart.ID {
			cartContent[i].Qty = cart.Qty
		}
	}
	json, err := json.Marshal(cartContent)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal mengubah data ke dalam format JSON.",
		})
		return
	}
	errSave := services.GetRedis().Set("cart" + userId, json, 0).Err()
	if errSave != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errSave.Error(),
			"result": nil,
			"description": "Tidak dapat menyimpan perubahan ke dalam keranjang belanja.",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": cartContent,
		"description": "Berhasil mengubah isi keranjang.",
	})
}