package user

import "golang.org/x/crypto/bcrypt"

type Service interface {
	CreateUser(input RegistrationInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateUser(input RegistrationInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Role = "user"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)

	newUser, err := s.repository.SaveUser(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}
