package repository

import (
	"bitka/services/auth/internal/domain"

	"gorm.io/gorm"
)

// TODO: Complete all methods of AuthRepository

type databaseRepo struct {
	db *gorm.DB
}

func NewDatabaseRepo(db *gorm.DB) domain.AuthRepository {
	return &databaseRepo{db: db}
}

func (r *databaseRepo) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *databaseRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *databaseRepo) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}
