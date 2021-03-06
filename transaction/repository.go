package transaction

import "gorm.io/gorm"

type Repository interface {
	FindTransactionByCampaignID(campaignID int) ([]Transaction, error)
	FindTransactionByUserID(userID int) ([]Transaction, error)
	FindTransactionByID(ID int) (Transaction, error)
	SaveTransaction(transaction Transaction) (Transaction, error)
	UpdateTransaction(transaction Transaction) (Transaction, error)
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

func (repo *repository) FindTransactionByUserID(userID int) ([]Transaction, error) {
	var transactionList []Transaction

	err := repo.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("created_at desc").Find(&transactionList).Error
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}

func (r *repository) FindTransactionByID(ID int) (Transaction, error) {
	var transaction Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repo *repository) SaveTransaction(transaction Transaction) (Transaction, error) {
	err := repo.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repo *repository) UpdateTransaction(transaction Transaction) (Transaction, error) {
	err := repo.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
