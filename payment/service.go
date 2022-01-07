package payment

import (
	"log"
	"os"
	"rocketship/campaign"
	"rocketship/user"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/veritrans/go-midtrans"
)

type service struct {
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewPaymentService(campaignRepository campaign.Repository) *service {
	return &service{campaignRepository}
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func (service *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = goDotEnvVariable("SERVER_KEY")
	midclient.ClientKey = goDotEnvVariable("CLIENT_KEY")
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
