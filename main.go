package main

import (
	"fmt"
	"log"
	"rocketship/auth"
	"rocketship/handler"
	"rocketship/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:1000)/rocketship?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewJWTService()
	userHandler := handler.NewUserHandler(userService, authService)

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOH0.nR_OBhtxFLbDYcqdsPckAIu9in0m1BH3QqnVK5HvXZY")
	if err != nil {
		fmt.Println(err.Error())
	}

	if token.Valid {
		fmt.Println("VALID")
	} else {
		fmt.Println("INVALID")
	}

	router := gin.Default()
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/validate_email", userHandler.ValidateEmail)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
}
