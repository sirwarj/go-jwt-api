package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirwarj/jwt-api/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterStru struct {
	Username string `json:"username" bindind:"required"`
	Password string `json:"password" bindind:"required"`
	Fullname string `json:"fullname" bindind:"required"`
	Avatar   string `json:"avatar" bindind:"required"`
}

type LoginStru struct {
	Username string `json:"username" bindind:"required"`
	Password string `json:"password" bindind:"required"`
}

var hmacSampleSecret []byte

func Register(ctx *gin.Context) {
	var json RegisterStru
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	// Check User Exists
	var userExist models.User
	models.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User Exists",
		})
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := models.User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname, Avatar: json.Avatar}
	models.Db.Save(&user)
	if user.ID > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "User Create Successful",
			"userId":  user.ID,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User Create Failed",
		})
	}
}

func Login(ctx *gin.Context) {
	var json LoginStru
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Check User Exists
	var userExist models.User
	models.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User Does Not Exists",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Login Failed",
		})
	} else {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Login Successful",
			"token":   tokenString,
		})
	}
}
