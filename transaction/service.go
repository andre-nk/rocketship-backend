package transaction

import (
	"errors"
	"rocketship/campaign"
	"rocketship/payment"
	"strconv"
)

type Service interface {
	FindTransactionByCampaignID(campaignID FindTransactionByIDInput) ([]Transaction, error)
	FindTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

type service struct {
	repository     Repository
	campaign       campaign.Repository
	paymentService payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (service *service) FindTransactionByCampaignID(input FindTransactionByIDInput) ([]Transaction, error) {
	campaign, err := service.campaign.FindCampaignByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("could not find transactions due to lack of credentials")
	}

	transactionList, err := service.repository.FindTransactionByCampaignID(input.ID)
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}

func (service *service) FindTransactionByUserID(userID int) ([]Transaction, error) {
	transactionList, err := service.repository.FindTransactionByUserID(userID)
	if err != nil {
		return transactionList, err
	}

	return transactionList, nil
}

func (service *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
		UserID:     input.User.ID,
		Status:     "pending",
	}

	newTransaction, err := service.repository.SaveTransaction(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := service.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	newTransaction, err = service.repository.UpdateTransaction(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (service *service) ProcessPayment(input TransactionNotificationInput) error {
	transactionID, _ := strconv.Atoi(input.OrderID)

	transaction, err := service.repository.FindTransactionByID(transactionID)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := service.repository.UpdateTransaction(transaction)
	if err != nil {
		return err
	}

	campaign, err := service.campaign.FindCampaignByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.FunderAmount = campaign.FunderAmount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := service.campaign.UpdateCampaign(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
