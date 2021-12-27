package user

import "gorm.io/gorm"

type Repository interface {
	SaveUser(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (repo *repository) SaveUser(user User) (User, error) {
	err := repo.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
