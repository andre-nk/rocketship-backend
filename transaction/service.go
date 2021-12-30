package transaction

import (
	"errors"
	"rocketship/campaign"
)

type Service interface {
	FindTransactionByCampaignID(campaignID TransactionByIDInput) ([]Transaction, error)
	FindTransactionByUserID(userID int) ([]Transaction, error)
}

type service struct {
	repository Repository
	campaign   campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (service *service) FindTransactionByCampaignID(input TransactionByIDInput) ([]Transaction, error) {
	campaign, err := service.campaign.FindCampaignByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.ID != input.User.ID {
		return []Transaction{}, errors.New("Could not find transactions due to lack of credentials")
	}

	transactionList, err := service.repository.FindTransactionByCampaignID(input.ID)
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}

func (service *service) FindTransactionByUserID(userID int) ([]Transaction, error) {
	// campaign, err := service.campaign.FindCampaignByID(input.ID)
	// if err != nil {
	// 	return []Transaction{}, err
	// }

	// if campaign.ID != input.User.ID {
	// 	return []Transaction{}, errors.New("Could not find transactions due to lack of credentials")
	// }

	transactionList, err := service.repository.FindTransactionByUserID(userID)
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}
