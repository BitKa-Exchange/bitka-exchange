package postgres

import (
	"bitka/services/auth/internal/domain"
	"errors"
	"gorm.io/gorm"
	"strings"
)

// TODO: Complete all methods of AuthRepository

type databaseRepo struct {
	db *gorm.DB
}

func NewDatabaseRepo(db *gorm.DB) domain.AuthRepository {
	return &databaseRepo{db: db}
}

func (r *databaseRepo) CreateUser(user *domain.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return errors.New("email already in use")
		}
		if strings.Contains(err.Error(), "users_username_key") {
			return errors.New("username already in use")
		}
		return err
	}
	return nil
}

func (r *databaseRepo) FindByEmailOrUser(identifier string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *databaseRepo) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}
