package handler

import (
	"net/http"
	"rocketship/campaign"
	"rocketship/helper"
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
