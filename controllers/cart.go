package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/twinj/uuid"
	"github.com/adeindriawan/itsfood-commerce/utils"
)

type Cart struct {
	ID string 					`json:"id"`
	MenuID uint64 			`json:"menu_id"`
	Name string					`json:"name"`
	Price uint64				`json:"price"`
	COGS uint64					`json:"cogs"`
	VendorID uint64			`json:"vendor_id"`
	Qty uint 						`json:"qty"`
	Image string				`json:"image"`
	MinOrderQty uint		`json:"min_order_qty"`
	MaxOrderQty uint		`json:"max_order_qty"`
	PreOrderHours uint	`json:"pre_order_hours"`
	PreOrderDays uint 	`json:"pre_order_days"`
}

func UserCartContent(userId string) ([]Cart, error) {
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

func _menuExistsAndChangeQty(menuId uint64, qty uint, cartContent []Cart) (bool, []Cart) {
	for i, search := range cartContent {
		if search.MenuID == menuId {
			cartContent[i].Qty += qty
			return true, cartContent
		}
	}
	return false, nil
}

func _menuMinOrderQtyValidated(menuQty uint, menuMinOrderQty uint) bool {
	return menuQty >= menuMinOrderQty
}

func _menuMaxOrderQtyValidated(menuQty uint, menuMaxOrderQty uint) bool {
	return menuQty <= menuMaxOrderQty
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

	user, err := utils.AuthCheck(c)
	if err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil token user yang ada. Unauthorized.",
		})
		return
	}

	if !_menuMinOrderQtyValidated(cart.Qty, cart.MinOrderQty) {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Menu min order qty failed to be validated.",
			"result": nil,
			"description": "Qty menu yang dimasukkan tidak sesuai dengan minimum order qty menu tersebut.",
		})

		return
	}

	if cart.MaxOrderQty != 0 && !_menuMaxOrderQtyValidated(cart.Qty, cart.MaxOrderQty) {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Menu max order qty failed to be validated.",
			"result": nil,
			"description": "Qty menu yang dimasukkan tidak sesuai dengan maksimum order qty menu tersebut.",
		})

		return
	}

	userId := strconv.Itoa(int(user))
	currentCartContent, noCurrentCartContent := UserCartContent(userId) // apa sudah ada cart dari user ini?
	if noCurrentCartContent != nil { // jika belum, maka buat cart baru untuk user ini
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
	user, err := utils.AuthCheck(c)
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
	cart, err := UserCartContent(userId)
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

	user, err := utils.AuthCheck(c)
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
	cartContent, err := UserCartContent(userId)
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
		if cart.ID == search.ID {
			if !_menuMinOrderQtyValidated(cart.Qty, search.MinOrderQty) {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": "Menu min order qty failed to be validated.",
					"result": nil,
					"description": "Qty menu yang dimasukkan tidak sesuai dengan minimum order qty menu tersebut.",
				})
				return
			} else if !_menuMaxOrderQtyValidated(cart.Qty, search.MaxOrderQty) {
				c.JSON(400, gin.H{
					"status": "failed",
					"errors": "Menu max order qty failed to be validated.",
					"result": nil,
					"description": "Qty menu yang dimasukkan tidak sesuai dengan maksimum order qty menu tersebut.",
				})
				return
			} else {
				cartContent[i].Qty = cart.Qty
			}
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

func DeleteCart(c *gin.Context) {
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

	user, err := utils.AuthCheck(c)
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
	cartContent, err := UserCartContent(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak ada isi keranjang belanja dari user yang bersangkutan.",
		})
		return
	}

	for i, v := range cartContent {
		if v.ID == cart.ID {
			cartContent[i] = cartContent[len(cartContent) - 1]
		}
	}

	newCartContent := cartContent[:len(cartContent) - 1]
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
			"description": "Tidak dapat mengeluarkan menu dari keranjang belanja.",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": newCartContent,
		"description": "Berhasil mengeluarkan menu dari keranjang belanja.",
	})
}

func CartTotals(c *gin.Context) {
	user, err := utils.AuthCheck(c)
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
	cartContent, err := UserCartContent(userId)
	
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak ada isi keranjang belanja dari user yang bersangkutan.",
		})
		return
	}

	totalQty := 0
	totalAmount := 0
	for _, v := range cartContent {
		totalQty += int(v.Qty)
		totalAmount += int(v.Price) * int(v.Qty)
	}

	type CartTotal struct {
		Items int		`json:"items"`
		Qty int			`json:"qty"`
		Amount int	`json:"amount"`
	}

	var total CartTotal

	total.Items = len(cartContent)
	total.Qty = totalQty
	total.Amount = totalAmount
	
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": total,
		"description": "Berhasil mengambil jumlah item di dalam keranjang belanja user.",
	})
}

func DestroyCart(c *gin.Context) {
	user, err := utils.AuthCheck(c)
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
	errDelete := services.GetRedis().Del("cart" + userId).Err()
	if errDelete != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errDelete.Error(),
			"result": nil,
			"description": "Gagal menghapus data keranjang belanja user.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": nil,
		"description": "Berhasil menghapus data keranjang belanja user.",
	})
}