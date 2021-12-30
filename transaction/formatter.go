package transaction

import "time"

type TransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	return formatter
}

func FormatTransactionList(transactionList []Transaction) []TransactionFormatter {
	var formatterList []TransactionFormatter

	for _, transaction := range transactionList {
		formatter := FormatTransaction(transaction)
		formatterList = append(formatterList, formatter)
	}

	return formatterList
}
