package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *models.User) error
	ReadUser(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
	ValidateCredentials(email, password string) (*models.User, error)
}

type UserServiceImpl struct {
	Repo repositories.UserRepository
}

func (s *UserServiceImpl) CreateUser(user *models.User) error {
	user.ID = uuid.Nil

	// By default, new users have the APP_USER role
	if user.Role == "" {
		user.Role = models.RoleAppUser
	}

	// Check if email already exists
	existingUser, _ := s.Repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)

	return s.Repo.CreateUser(user)
}

func (s *UserServiceImpl) ReadUser(id uuid.UUID) (*models.User, error) {
	return s.Repo.ReadUser(id)
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	// Check if user exists
	existingUser, err := s.Repo.ReadUser(user.ID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// If a new password is provided, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)
	} else {
		// Keep of one
		user.Password = existingUser.Password
	}

	return s.Repo.UpdateUserSafe(user)
}

func (s *UserServiceImpl) DeleteUser(id uuid.UUID) error {
	// Check if user exists
	if _, err := s.Repo.ReadUser(id); err != nil {
		return fmt.Errorf("user not found")
	}
	return s.Repo.DeleteUser(id)
}

func (s *UserServiceImpl) ValidateCredentials(email, password string) (*models.User, error) {
	// Fetch user by email
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Ne pas renvoyer le hash du mot de passe
	user.Password = ""

	return user, nil
}
