package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(input RegistrationInput) (User, error)
	Login(input LoginInput) (User, error)
	ValidateEmail(email EmailValidatorInput) (bool, error)
	UploadAvatar(id int, filePath string) (User, error)
	FindUserByID(id int) (User, error)
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

	newUser, err := s.repository.CreateUser(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindUserByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that e-mail")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) ValidateEmail(validatorInput EmailValidatorInput) (bool, error) {
	user, err := s.repository.FindUserByEmail(validatorInput.Email)
	if err != nil {
		return false, err
	}

	if user.ID != 0 {
		return false, err
	}

	return true, nil
}

func (s *service) UploadAvatar(id int, filePath string) (User, error) {
	user, err := s.repository.FindUserByID(id)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = filePath

	updatedUser, err := s.repository.UpdateUser(user)
	if err != nil {
		return user, err
	}

	return updatedUser, nil
}

func (s *service) FindUserByID(id int) (User, error) {
	user, err := s.repository.FindUserByID(id)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with this ID")
	}

	return user, nil
}
