package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	FindCampaigns(userID int) ([]Campaign, error)
	FindCampaignByID(campaignID CampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(campaignID CampaignDetailInput, input CreateCampaignInput) (Campaign, error)
	CreateCampaignImage(input CreateCampaignImageInput, filePath string) (CampaignImage, error)
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

func (s *service) UpdateCampaign(campaignID CampaignDetailInput, input CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindCampaignByID(campaignID.ID)
	if err != nil {
		return campaign, err
	}

	if input.User.ID != campaign.User.ID {
		return campaign, errors.New("could not update this campaign due to lack of credentials")
	}

	campaign.Name = input.Name
	campaign.Description = input.Description
	campaign.ShortDescription = input.ShortDescription
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	slugWireframe := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugWireframe)

	updatedCampaign, err := s.repository.UpdateCampaign(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) CreateCampaignImage(input CreateCampaignImageInput, filePath string) (CampaignImage, error) {
	campaign, err := s.repository.FindCampaignByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.ID != input.User.ID {
		return CampaignImage{}, errors.New("could not upload campaign image due to lack of credentials")
	}

	isPrimary := 0

	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{
		FileName:   filePath,
		CampaignID: input.CampaignID,
		IsPrimary:  isPrimary,
	}

	createdImage, err := s.repository.UploadCampaignImage(campaignImage)
	if err != nil {
		return createdImage, err
	}

	return createdImage, nil
}
