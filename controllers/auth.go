package controllers

import (
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/adeindriawan/itsfood-commerce/models"
	jwt "github.com/golang-jwt/jwt/v4"
	redis "github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)

var  client *redis.Client

func init() {
	//Initializing redis
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  dsn := os.Getenv("REDIS_HOST")
  if len(dsn) == 0 {
     dsn = "localhost:6379"
  }
	fmt.Println(dsn)
  client = redis.NewClient(&redis.Options{
     Addr: dsn, //redis port
  })
  _, errRedis := client.Ping().Result()
  if errRedis != nil {
     panic(err)
  }
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
	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": register})
}

type UserLoginInput struct {
	Email string 		`json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var user models.User
	var login UserLoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		c.AbortWithStatus(401)
		fmt.Println(err)
	} else if user.Email != login.Email || user.Password != login.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to authenticate"})
	} else {
		ts, err := CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		saveErr := CreateAuth(user.ID, ts)
		if saveErr != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
		}
		data := map[string]interface{}{
			"user": user,
			"token": ts,
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"errors": nil,
			"result": data,
			"description": "Berhasil login",
		})
	}
}

type TokenDetails struct {
  AccessToken  string
  RefreshToken string
  AccessUuid   string
  RefreshUuid  string
  AtExpires    int64
  RtExpires    int64
}

func CreateToken(userId uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
  td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
  td.AccessUuid = uuid.NewV4().String()

  td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
  td.RefreshUuid = uuid.NewV4().String()

	var err error
	// creating access token
	os.Setenv("ACCESS_SECRET", "loremipsum")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// creating refresh token
	// os.Setenv("REFRESH_SECRET", "loremipsum")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateAuth(userId uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC (to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title string `json:"title"`
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

type AccessDetails struct {
	AccessUuid string
	UserId uint64
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId: userId,
		}, nil
	}

	return nil, err
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
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
	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak ada token user yang sesuai. Unauthorized.",
		})
		return
	}
	userId, err := FetchAuth(tokenAuth)
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
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func Logout(c *gin.Context) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"errors": err.Error(),
			"result": nil,
			"description": "Tidak dapat mengekstrak token user.",
		})
		return
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

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"errors": err.Error(),
				"result": nil,
				"description": "Token dari user tidak valid.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
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
		ts, createErr := CreateToken(userId)
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
		saveErr := CreateAuth(userId, ts)
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