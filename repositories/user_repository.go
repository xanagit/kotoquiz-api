package repositories

import (
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	ReadUser(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) ReadUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUser(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, "id = ?", id).Error
}
