package repositories

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	ReadUser(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserSafe(user *models.User) error
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

func (r *UserRepositoryImpl) UpdateUserSafe(user *models.User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// SERIALIZABLE ensures that no other transactions can modify
		// the data while we are modifying it
		err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").Error
		if err != nil {
			return err
		}

		// Lock the record before modification
		var existingUser models.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			First(&existingUser, "id = ?", user.ID).Error; err != nil {
			return err
		}

		// Check if the email is already used by another user
		if user.Email != existingUser.Email {
			var count int64
			if err := tx.Model(&models.User{}).
				Where("email = ? AND id != ?", user.Email, user.ID).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return fmt.Errorf("email already in use")
			}
		}

		return tx.Save(user).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (r *UserRepositoryImpl) DeleteUser(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, "id = ?", id).Error
}
