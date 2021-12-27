package user

import "gorm.io/gorm"

type Repository interface {
	CreateUser(user User) (User, error)
	FindUserByEmail(email string) (User, error)
	FindUserByID(id int) (User, error)
	UpdateUser(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (repo *repository) CreateUser(user User) (User, error) {
	err := repo.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *repository) FindUserByEmail(email string) (User, error) {
	var user User
	err := repo.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *repository) FindUserByID(id int) (User, error) {
	var user User
	err := repo.db.Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *repository) UpdateUser(user User) (User, error) {
	err := repo.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
