package controllers

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adeindriawan/itsfood-commerce/models"
	"github.com/adeindriawan/itsfood-commerce/services"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/twinj/uuid"
	"github.com/adeindriawan/itsfood-commerce/utils"
)

type ForgotPasswordPayload struct {
	Email string `json:"email"`
}

func ForgotPassword(c *gin.Context) {
	var payload ForgotPasswordPayload
	var user models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}
	query := services.DB.First(&user, "email = ?", payload.Email)
	if query.Error != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": query.Error.Error(),
			"result": nil,
			"description": "Gagal melakukan query.",
		})
		return
	}

	resetToken := uuid.NewV4().String()
	mailTo := user.Email
	mailSubject := "[ITS Food] Lupa Kata Sandi"
	mailBody := resetToken

	resetTokenExpires := time.Now().Add(time.Minute * 15).Unix()
	rtx := time.Unix(resetTokenExpires, 0)
	now := time.Now()
	if err := services.GetRedis().Set(resetToken, mailTo, rtx.Sub(now)).Err(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal menyimpan reset token di sistem.",
		})
		return
	}

	_, errorSendingResetPasswordEmail := services.SendMail(mailTo, mailSubject, mailBody)
	if errorSendingResetPasswordEmail != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": nil,
			"result": nil,
			"description": "Gagal mengirim email berisi token ke alamat " + mailTo,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": user,
		"description": "Sukses mengirim email berisi token ke alamat " + mailTo,
	})
}

type ResetPasswordPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ResetPassword(c *gin.Context) {
	var payload ResetPasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Gagal memvalidasi data yang masuk",
			"result": nil,
			"description": "Data password tidak sama dengan confirm password yang dikirim.",
		})
		return
	}

	if email, err := services.GetRedis().Get(payload.Token).Result(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Token tidak ditemukan dalam sistem. Kemungkinan sudah kadaluwarsa.",
		})
		return
	} else if email != payload.Email {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Email yang terkirim tidak sama dengan email yang tersimpan dalam token di Redis.",
			"result": nil,
			"description": "User dengan email ini tidak dapat mereset password.",
		})
		return
	}

	var user models.User
	findUser := services.DB.First(&user, "email = ?", payload.Email)
	if findUser.Error != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": findUser.Error.Error(),
			"result": nil,
			"description": "Gagal menemukan user dengan email tersebut dalam sistem.",
		})
		return
	}

	hash, errHash := utils.HashPassword(payload.Password)
	fmt.Println(hash)
	if errHash != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errHash.Error(),
			"result": nil,
			"descripion": "Gagal membuat hash dari password yang diberikan.",
		})
		return
	}
	user.Password = hash
	fmt.Println(hash)
	updatePassword := services.DB.Save(&user)
	if updatePassword.Error != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": updatePassword.Error.Error(),
			"result": nil,
			"description": "Gagal mengubah password dari user ini.",
		})
		return
	}
	fmt.Println(hash)
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": hash,
		"description": "Sukses mengganti password dari user ini.",
	})
}

type AdminRegisterInput struct {
	Name string				`json:"name"`
	Email string			`json:"email"`
	Password string 	`json:"password"`
	Phone string 			`json:"phone"`
}

func AdminRegister(c *gin.Context) {
	var register AdminRegisterInput
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	hashedPassword, errHashingPassword := utils.HashPassword(register.Password)
	if errHashingPassword != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errHashingPassword.Error(),
			"result": nil,
			"description": "Gagal membuat hash password.",
		})
		return
	}

	user := models.User{Name: register.Name, Email: register.Email, Password: hashedPassword, Phone: register.Phone, Type: "Admin", Status: "Registered", CreatedBy: register.Name, UpdatedAt: time.Time{}}
	if errorCreatingUser := services.DB.Create(&user).Error; errorCreatingUser != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingUser.Error(),
			"result": nil,
			"description": "Gagal menyimpan data user baru dalam database.",
		})
		return
	}

	userId := user.ID
	admin := models.Admin{UserID: userId, Name: register.Name, Email: register.Email, Phone: register.Phone, Status: "Inactive", CreatedBy: register.Name, UpdatedAt: time.Time{}}
	if errorCreatingAdmin := services.DB.Create(&admin).Error; errorCreatingAdmin != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingAdmin.Error(),
			"result": nil,
			"description": "Gagal menyimpan data admin baru dalam database.",
		})
		return
	}
	
	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": user,
		"description": "Berhasil menambah admin baru.",
	})
}

type CustomerRegisterInput struct {
	Name string			`json:"name"`
	Email string		`json:"email"`
	Password string `json:"password"`
	Phone string 		`json:"phone"`
	Type string 		`json:"type"`
	UnitID uint64 	`json:"unit_id"`
}

func CustomerRegister(c *gin.Context) {
	var register CustomerRegisterInput
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk",
		})
		return
	}

	hashedPassword, errorHashingPassword := utils.HashPassword(register.Password)
	if errorHashingPassword != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorHashingPassword.Error(),
			"result": nil,
			"description": "Gagal membuat hash password.",
		})
		return
	}

	user := models.User{Name: register.Name, Email: register.Email, Password: hashedPassword, Phone: register.Phone, Type: "Customer", Status: "Registered", CreatedBy: register.Name, UpdatedAt: time.Time{}}
	if errorCreatingUser := services.DB.Create(&user).Error; errorCreatingUser != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingUser.Error(),
			"result": nil,
			"description": "Gagal menyimpan data user baru dalam database.",
		})
		return
	}

	userId := user.ID
	customer := models.Customer{UserID: userId, Type: register.Type, UnitID: register.UnitID, Status: "Registered", CreatedBy: register.Name, UpdatedAt: time.Time{}}
	if errorCreatingCustomer := services.DB.Create(&customer).Error; errorCreatingCustomer != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingCustomer.Error(),
			"result": nil,
			"description": "Gagal menyimpan data customer baru dalam database.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": nil,
		"description": "Berhasil mendaftarkan customer baru ke dalam sistem.",
	})
}

type VendorRegisterInput struct {
	Name string 							`json:"name"`
	Email string 							`json:"email"`
	Password string						`json:"password"`
	Phone string 							`json:"phone"`
	CompanyName string				`json:"company_name"`
	CompanyType string 				`json:"company_type"`
	Address string 						`json:"address"`
	BankName string 					`json:"bank_name"`
	BankBranch string 				`json:"bank_branch"`
	BankAccountNumber	string	`json:"bank_account_number"`
	BankAccountName string 		`json:"bank_account_name"`
}

func VendorRegister(c *gin.Context) {
	var register VendorRegisterInput
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	hashedPassword, errorHashingPassword := utils.HashPassword(register.Password)
	if errorHashingPassword != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorHashingPassword.Error(),
			"result": nil,
			"description": "Gagal membuat hash password.",
		})
		return
	}

	user := models.User{Name: register.Name, Email: register.Email, Password: hashedPassword, Phone: register.Phone, Status: "Registered", Type: "Vendor", CreatedBy: register.Name, CreatedAt: time.Time{}}
	if errorCreatingUser := services.DB.Create(&user).Error; errorCreatingUser != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingUser.Error(),
			"result": nil,
			"description": "Gagal mennyimpan data user baru.",
		})
		return
	}

	userId := user.ID
	vendor := models.Vendor{
		UserID: userId,
		CompanyName: register.CompanyName, 
		CompanyType: register.CompanyType,
		Address: register.Address,
		BankName: register.BankName,
		BankBranch: register.BankBranch,
		BankAccountNumber: register.BankAccountNumber,
		BankAccountName: register.BankAccountName,
	}
	if errorCreatingVendor := services.DB.Create(&vendor).Error; errorCreatingVendor != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": errorCreatingVendor.Error(),
			"result": nil,
			"description": "Gagal menyimpan data vendor baru.",
		})
		return
	}
	data := map[string]interface{}{
		"user": user,
		"vendor": vendor,
	}

	c.JSON(200, gin.H{
		"status": "success",
		"errors": nil,
		"result": data,
		"description": "Berhasil mendaftarkan vendor baru.",
	})
}

type UserRegisterInput struct {
	Email string 		`json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var register UserRegisterInput
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Email: register.Email, Password: register.Password}
	services.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": register})
}

type UserLoginInput struct {
	Email string 		`json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CustomerLogin(c *gin.Context) {
	var user models.User
	var login UserLoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	if err := services.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal menemukan user dengan email yang dikirimkan.",
		})
		return
	} else if user.Email != login.Email || !utils.CheckPasswordHash(login.Password, user.Password) {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Gagal mengautentikasi.",
			"result": nil,
			"description": "Gagal mengautentikasi info login dari data yang dikirimkan.",
		})
		return
	} else {
		if user.Type != "Customer" {
			c.JSON(401, gin.H{
				"status": "failed",
				"errors": "Not customer",
				"result": nil,
				"description": "User yang bersangkutan bukan bertipe Customer.",
			})
			return
		}
		var customer models.Customer
		if err := services.DB.Preload("User").Where("user_id = ?", user.ID).First(&customer).Error; err != nil {
			c.JSON(401, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Gagal menemukan data customer dengan ID user tersebut.",
			})
			return
		}
		ts, err := utils.CreateToken(user.ID)
		if err != nil {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Tidak dapat membuat token untuk proses autentikasi.",
			})
			return
		}
		saveErr := utils.CreateAuth(user.ID, ts)
		if saveErr != nil {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": saveErr.Error(),
				"result": nil,
				"description": "Gagal membuat autentikasi user.",
			})
			return
		}
		data := map[string]interface{}{
			"token": ts,
			"customer": customer,
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"errors": nil,
			"result": data,
			"description": "Berhasil login",
		})
	}
}

func AdminLogin(c *gin.Context) {
	var user models.User
	var login UserLoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal memproses data yang masuk.",
		})
		return
	}

	if err := services.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal menemukan user dengan email yang dikirimkan.",
		})
		return
	} else if user.Email != login.Email || !utils.CheckPasswordHash(login.Password, user.Password) {
		c.JSON(400, gin.H{
			"status": "failed",
			"errors": "Gagal mengautentikasi.",
			"result": nil,
			"description": "Gagal mengautentikasi info login dari data yang dikirimkan.",
		})
		return
	} else {
		if user.Type != "Admin" {
			c.JSON(401, gin.H{
				"status": "failed",
				"errors": "Not admin",
				"result": nil,
				"description": "User yang bersangkutan bukan bertipe Admin.",
			})
			return
		}
		var admin models.Admin
		if err := services.DB.Preload("User").Where("user_id = ?", user.ID).First(&admin).Error; err != nil {
			c.JSON(401, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Gagal menemukan data admin dengan ID user tersebut.",
			})
			return
		}
		ts, err := utils.CreateToken(user.ID)
		if err != nil {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Tidak dapat membuat token untuk proses autentikasi.",
			})
			return
		}
		saveErr := utils.CreateAuth(user.ID, ts)
		if saveErr != nil {
			c.JSON(422, gin.H{
				"status": "failed",
				"errors": saveErr.Error(),
				"result": nil,
				"description": "Gagal membuat autentikasi user.",
			})
			return
		}
		data := map[string]interface{}{
			"token": ts,
			"admin": admin,
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"errors": nil,
			"result": data,
			"description": "Berhasil login",
		})
	}
}

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title string `json:"title"`
}

func CreateTodo(c *gin.Context) {
	var td *Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat memproses data yang masuk. Invalid JSON",
		})
		return
	}
	tokenAuth, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak ada token user yang sesuai. Unauthorized.",
		})
		return
	}
	userId, err := utils.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengambil token user yang ada. Unauthorized.",
		})
		return
	}
	td.UserID = userId

	// you can proceed to save the Todo to a database
	// but we will just return it to the caller here
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"errors": nil,
		"result": td,
		"description": "Berhasil membuat data baru.",
	})
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := services.GetRedis().Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func Logout(c *gin.Context) {
	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengekstrak token user.",
		})
		return
	}
	
	_, errorFetchingAuth := utils.FetchAuth(au)
	if errorFetchingAuth != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal mengambil user ID.",
		})
		return
	}

	if checkCustomerContext, exists := c.Get("customer"); exists {
		customerContext := checkCustomerContext.(models.Customer)
		customerId := strconv.Itoa(int(customerContext.ID))
		errDestroy := DestroyCustomerCart(customerId)
		if errDestroy != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"errors": errDestroy.Error(),
				"result": nil,
				"description": "Gagal menghapus data keranjang belanja user.",
			})
			return
		}	
	}

	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": "Tidak ada token user yang terhapus: " + delErr.Error(),
			"result": nil,
			"description": "Error dalam menghapus token user atau tidak ada token yang terhapus.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"errors": nil,
		"result": nil,
		"description": "Berhasil log out.",
	})
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Data yang masuk tidak dapat diproses lebih lanjut.",
		})
		return
	}
	refreshSecret := os.Getenv("REFRESH_SECRET")
	refreshToken := mapToken["refresh_token"]
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected string method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})
	// If there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Refresh token yang ada sudah kadaluarsa. Silakan login kembali.",
		})
		return
	}
	// Is the token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Gagal membuat access & refresh token yang baru.",
		})
		return
	}
	// Since the token is valid, get the uuid
	claims, ok := token.Claims.(jwt.MapClaims) // the token should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) // convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Gagal membuat access & refresh token yang baru.",
			})
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Tidak bisa membuat access & refresh token yang baru karena gagal mengonversi ID user.",
			})
			return
		}
		// Delete the previous refresh token
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"errors": delErr.Error(),
				"result": nil,
				"description": "Tidak bisa membuat access & refresh token yang baru karena gagal menghapus refresh token yang lama.",
			})
			return
		}
		// Create new pairs of refresh and access token
		ts, createErr := utils.CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "failed",
				"errors": createErr.Error(),
				"result": nil,
				"description": "Gagal membuat access & refresh token yang baru.",
			})
			return
		}
		// Save the token metadata to Redis
		saveErr := utils.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "failed",
				"errors": saveErr.Error(),
				"result": nil,
				"description": "Gagal menyimpan metadata token ke Redis.",
			})
			return
		}
		tokens := map[string]string{
			"access_token": ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"errors": nil,
			"result": tokens,
			"description": "Berhasil memperbarui access & refresh token.",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": "Refresh expired",
			"result": nil,
			"description": "Token untuk merefresh access token sudah kadaluarsa.",
		})
	}
}
