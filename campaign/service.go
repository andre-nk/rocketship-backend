package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	FindCampaigns(userID int) ([]Campaign, error)
	FindCampaignByID(campaignID CampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindCampaignByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAllCampaign()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) FindCampaignByID(campaignID CampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindCampaignByID(campaignID.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		UserID:           input.User.ID,
	}

	slugWireframe := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugWireframe)

	newCampaign, err := s.repository.CreateCampaign(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}
