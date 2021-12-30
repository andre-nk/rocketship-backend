package handler

import (
	"fmt"
	"net/http"
	"rocketship/auth"
	"rocketship/helper"
	"rocketship/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (handler *userHandler) RegisterUser(context *gin.Context) {
	var input user.RegistrationInput

	err := context.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Account registration failed due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			errorMessage,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := handler.userService.CreateUser(input)
	if err != nil {
		response := helper.APIResponse(
			"Account registration failed due to server error",
			http.StatusBadRequest,
			"failed",
			nil,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := handler.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse(
			"Account registration failed due to token generation error",
			http.StatusBadRequest,
			"failed",
			nil,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	formattedUser := user.FormatUser(newUser, token)
	response := helper.APIResponse(
		"Account has been registered",
		http.StatusOK,
		"success",
		formattedUser,
	)

	context.JSON(http.StatusOK, response)
}

func (handler *userHandler) Login(context *gin.Context) {
	var input user.LoginInput

	err := context.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Log in failed due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			errorMessage,
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := handler.userService.Login(input)

	if err != nil {
		response := helper.APIResponse(
			"Log in failed due to wrong credentials",
			http.StatusBadRequest,
			"failed",
			nil,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if loggedUser.ID == 0 {
		response := helper.APIResponse(
			"There is no user associated with this e-mail address",
			http.StatusBadRequest,
			"failed",
			nil,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := handler.authService.GenerateToken(loggedUser.ID)

	if err != nil {
		response := helper.APIResponse(
			"Log in failed due to token generation error",
			http.StatusBadRequest,
			"failed",
			nil,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	formattedUser := user.FormatUser(loggedUser, token)
	response := helper.APIResponse(
		"Log in successful",
		http.StatusOK,
		"success",
		formattedUser,
	)

	context.JSON(http.StatusOK, response)
}

func (handler *userHandler) ValidateEmail(context *gin.Context) {
	var input user.EmailValidatorInput

	err := context.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"E-mail validation failed due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			errorMessage,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := handler.userService.ValidateEmail(input)
	if err != nil {
		errors := helper.FormatError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"E-mail validation failed due to server error",
			http.StatusBadGateway,
			"failed",
			errorMessage,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "This e-mail address is available"
	if !isEmailAvailable {
		metaMessage = "This e-mail address is used"
	}

	response := helper.APIResponse(
		metaMessage,
		http.StatusOK,
		"success",
		data,
	)

	context.JSON(http.StatusOK, response)
}

func (handler *userHandler) UploadAvatar(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	file, err := context.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse(
			"Failed to upload avatar due to bad input",
			http.StatusUnprocessableEntity,
			"error",
			data,
		)

		context.JSON(http.StatusUnprocessableEntity, response)
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = context.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse(
			"Failed to save avatar due to server error",
			http.StatusBadRequest,
			"error",
			data,
		)

		context.JSON(http.StatusBadRequest, response)
	}

	_, err = handler.userService.UploadAvatar(userID, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse(
			"Failed to upload avatar due to server error",
			http.StatusBadRequest,
			"error",
			data,
		)

		context.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse(
		"User's avatar updated",
		http.StatusOK,
		"success",
		data,
	)

	context.JSON(http.StatusOK, response)
}
