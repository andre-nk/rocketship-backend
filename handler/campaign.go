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

func (handler *campaignHandler) FindCampaign(context *gin.Context) {
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
