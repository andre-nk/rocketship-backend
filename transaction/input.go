package transaction

import "rocketship/user"

type TransactionByIDInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
