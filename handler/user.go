package handler

import (
	"net/http"
	"rocketship/helper"
	"rocketship/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	formattedUser := user.FormatUser(newUser, "mockToken")
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
		context.JSON(http.StatusBadRequest, response)
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

	formattedUser := user.FormatUser(loggedUser, "mockToken")
	response := helper.APIResponse(
		"Log in successful",
		http.StatusOK,
		"success",
		formattedUser,
	)

	context.JSON(http.StatusOK, response)
}
