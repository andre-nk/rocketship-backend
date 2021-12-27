package user

import "time"

type User struct {
	ID             int
	Name           string
	Role           string
	Email          string
	PasswordHash   string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
