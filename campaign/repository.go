package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAllCampaign() ([]Campaign, error)
	FindCampaignByUserID(userID int) ([]Campaign, error)
	FindCampaignByID(campaignID int) (Campaign, error)
	CreateCampaign(campaign Campaign) (Campaign, error)
	UpdateCampaign(campaign Campaign) (Campaign, error)
	UploadCampaignImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllAsNonPrimary(campaignID int) (bool, error)
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

func (repo *repository) UpdateCampaign(campaign Campaign) (Campaign, error) {
	err := repo.db.Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (repo *repository) UploadCampaignImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := repo.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (repo *repository) MarkAllAsNonPrimary(campaignID int) (bool, error) {
	err := repo.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
