package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAllCampaign() ([]Campaign, error)
	FindCampaignByUserID(userID int) ([]Campaign, error)
	FindCampaignByID(campaignID int) (Campaign, error)
	CreateCampaign(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (repo *repository) FindAllCampaign() ([]Campaign, error) {
	var campaignList []Campaign

	err := repo.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaignList).Error
	if err != nil {
		return campaignList, err
	}

	return campaignList, nil
}

func (repo *repository) FindCampaignByUserID(userID int) ([]Campaign, error) {
	var campaignList []Campaign
	err := repo.db.Where("user_id =?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaignList).Error

	if err != nil {
		return campaignList, err
	}

	return campaignList, nil
}

func (repo *repository) FindCampaignByID(campaignID int) (Campaign, error) {
	var campaign Campaign

	err := repo.db.Preload("User").Preload("CampaignImages").Where("id = ?", campaignID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (repo *repository) CreateCampaign(campaign Campaign) (Campaign, error) {
	err := repo.db.Create(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
