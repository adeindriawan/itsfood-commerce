package utils

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"errors"
	"net/http"
	"github.com/twinj/uuid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/adeindriawan/itsfood-commerce/services"
	"github.com/adeindriawan/itsfood-commerce/models"
)

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
  td.AtExpires = time.Now().Add(time.Minute * 15).UnixMilli()
  td.AccessUuid = uuid.NewV4().String()

  td.RtExpires = time.Now().Add(time.Hour * 24 * 7).UnixMilli()
  td.RefreshUuid = uuid.NewV4().String()

	var err error
	// creating access token
	os.Setenv("ACCESS_SECRET", "loremipsum")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).UnixMilli()
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

	errAccess := services.GetRedis().Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := services.GetRedis().Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
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
	userid, err := services.GetRedis().Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func GetCustomerType(c *gin.Context) (string, error) {
	var errorCustomerNotIdentified error = errors.New("customer tidak teridentifikasi: customer belum login atau register")
	var customer models.Customer

	userId, err := AuthCheck(c)
	if err != nil {
		return "", errorCustomerNotIdentified
	}

	query := services.DB.First(&customer, userId)
	if query.Error != nil {
		return "", errorCustomerNotIdentified
	}

	if query.RowsAffected == 1 {
		customerType := customer.Type
		return customerType, nil
	}

	return "", errorCustomerNotIdentified
}