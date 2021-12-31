package transaction

import "time"

type TransactionFormatter struct {
	ID         int    `json:"id"`
	Amount     int    `json:"amount"`
	UserID     int    `json:"user_id"`
	CampaignID int    `json:"campaign_id"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionFormatter struct {
	ID        int                   `json:"id"`
	Amount    int                   `json:"amount"`
	Status    string                `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	Campaign  CampaignInfoFormatter `json:"campaign"`
}

type CampaignInfoFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserID:     transaction.UserID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
	}

	return formatter
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	return formatter
}

func FormatCampaignTransactionList(transactionList []Transaction) []CampaignTransactionFormatter {
	var formatterList []CampaignTransactionFormatter

	for _, transaction := range transactionList {
		formatter := FormatCampaignTransaction(transaction)
		formatterList = append(formatterList, formatter)
	}

	return formatterList
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
	}

	campaignFormatter := CampaignInfoFormatter{
		Name: transaction.Campaign.Name,
	}
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

func FormatUserTransactionList(transactionList []Transaction) []UserTransactionFormatter {
	var formatterList []UserTransactionFormatter

	for _, transaction := range transactionList {
		formatter := FormatUserTransaction(transaction)
		formatterList = append(formatterList, formatter)
	}

	return formatterList
}
