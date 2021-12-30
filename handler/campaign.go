package handler

import (
	"fmt"
	"net/http"
	"rocketship/campaign"
	"rocketship/helper"
	"rocketship/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (handler *campaignHandler) FindCampaigns(context *gin.Context) {
	userID, _ := strconv.Atoi(context.Query("user_id"))

	campaigns, err := handler.service.FindCampaigns(userID)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get campaigns due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse(
		"Campaigns fetched!",
		http.StatusOK,
		"success",
		campaign.FormatCampaigns(campaigns),
	)

	context.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) FindCampaign(context *gin.Context) {
	var input campaign.CampaignDetailInput

	err := context.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get campaign with that ID",
			http.StatusUnprocessableEntity,
			"error",
			err,
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignByID, err := handler.service.FindCampaignByID(input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get campaign with that ID due to server error",
			http.StatusBadRequest,
			"error",
			err,
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Campaign fetched",
		http.StatusOK,
		"success",
		campaign.FormatCampaignDetail(campaignByID),
	)
	context.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) CreateCampaign(context *gin.Context) {
	var input campaign.CreateCampaignInput

	err := context.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to create campaign due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := context.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := handler.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to create campaign due to server error",
			http.StatusBadRequest,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Campaign successfully created!",
		http.StatusOK,
		"success",
		campaign.FormatCampaign(newCampaign),
	)
	context.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) UpdateCampaign(context *gin.Context) {
	var inputID campaign.CampaignDetailInput
	var input campaign.CreateCampaignInput

	err := context.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse(
			"Failed to update campaign with that ID",
			http.StatusUnprocessableEntity,
			"error",
			err,
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = context.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to update campaign due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := context.MustGet("currentUser").(user.User)
	input.User = currentUser

	updatedCampaign, err := handler.service.UpdateCampaign(inputID, input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to update campaign due to server error",
			http.StatusBadRequest,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Campaign successfully updated!",
		http.StatusOK,
		"success",
		campaign.FormatCampaign(updatedCampaign),
	)
	context.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) UploadCampaignImage(context *gin.Context) {
	var input campaign.CreateCampaignImageInput
	currentUser := context.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	err := context.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to upload campaign image due to bad input",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := context.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse(
			"Failed to upload campaign image due to bad input",
			http.StatusUnprocessableEntity,
			"error",
			data,
		)

		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = currentUser

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = context.SaveUploadedFile(file, path)

	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		data := gin.H{
			"errors": errorMessage,
		}

		response := helper.APIResponse(
			"Failed to save campaign image due to server error",
			http.StatusBadRequest,
			"error",
			data,
		)

		context.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = handler.service.CreateCampaignImage(input, path)
	if err != nil {
		response := helper.APIResponse(
			"Failed to upload campaign image due to server error",
			http.StatusBadRequest,
			"error",
			err.Error(),
		)

		context.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse(
		"Campaign image uploaded",
		http.StatusOK,
		"success",
		data,
	)

	context.JSON(http.StatusOK, response)
}
