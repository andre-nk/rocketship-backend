package payment

import (
	"rocketship/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct{}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewPaymentService() *service {
	return &service{}
}

func (service *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResponse, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResponse.RedirectURL, nil
}
