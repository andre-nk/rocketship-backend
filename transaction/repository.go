package transaction

import "gorm.io/gorm"

type Repository interface {
	FindTransactionByCampaignID(campaignID int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (repo *repository) FindTransactionByCampaignID(campaignID int) ([]Transaction, error) {
	var transactionList []Transaction

	err := repo.db.Preload("User").Where("campaign_id = ?", campaignID).Order("created_at desc").Find(&transactionList).Error
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}
