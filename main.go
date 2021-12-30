package main

import (
	"log"
	"net/http"
	"rocketship/auth"
	"rocketship/campaign"
	"rocketship/handler"
	"rocketship/helper"
	"rocketship/transaction"
	"rocketship/user"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	//AUTH
	authService := auth.NewJWTService()

	//USER
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)

	//CAMPAIGN
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	//TRANSACTION
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//SANDBOX HERE===========================================

	//SANDBOX END============================================

	//ROUTER CONFIG
	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("api/v1")

	//AUTH ROUTES
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/validate_email", userHandler.ValidateEmail)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	//CAMPAIGN ROUTES
	api.GET("/campaigns", campaignHandler.FindCampaigns)
	api.GET("/campaigns/:id", campaignHandler.FindCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadCampaignImage)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)

	//TRANSACTION ROUTES
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.FindTransactionByCampaignID)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.FindTransactionByUserID)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(
				"Unauthorized request",
				http.StatusUnauthorized,
				"",
				nil,
			)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		tokenArray := strings.Split(authHeader, " ")
		if len(tokenArray) == 2 {
			tokenString = tokenArray[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse(
				"Unauthorized request due to invalid token",
				http.StatusUnauthorized,
				"",
				nil,
			)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse(
				"Unauthorized request due to invalid token",
				http.StatusUnauthorized,
				"",
				nil,
			)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.FindUserByID(userID)

		if err != nil {
			response := helper.APIResponse(
				"User search failed due to server error",
				http.StatusBadRequest,
				"",
				nil,
			)
			context.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		context.Set("currentUser", user)
	}
}
