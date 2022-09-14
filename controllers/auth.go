package controllers

import (
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
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
		c.JSON(http.StatusOK, gin.H{"data": data})
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
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// creating refresh token
	os.Setenv("REFRESH_SECRET", "loremipsum")
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
