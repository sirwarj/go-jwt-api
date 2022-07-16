package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirwarj/jwt-api/controllers/auth"
	"github.com/sirwarj/jwt-api/controllers/middleware"
	"github.com/sirwarj/jwt-api/controllers/user"
	"github.com/sirwarj/jwt-api/models"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	models.InitDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	authorized := r.Group("/users", middleware.JWTAuth())
	authorized.GET("/readall", user.ReadAll)
	authorized.GET("/profile", user.Profile)

	r.Run("localhost:8080")
}
