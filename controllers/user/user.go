package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirwarj/jwt-api/models"
)

var hmacSampleSecret []byte

func ReadAll(ctx *gin.Context) {
	var users []models.User
	models.Db.Find(&users)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User Read Success",
		"users":   users,
	})
}

func Profile(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(float64)
	var user models.User
	models.Db.First(&user, userId)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})

}
