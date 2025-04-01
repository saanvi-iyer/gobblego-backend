package user

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(db *gorm.DB, user *models.User) error
	GetUserByID(db *gorm.DB, id string) (*models.User, error)
	GetUsersByCartID(db *gorm.DB, cartID string) ([]models.User, error)
	SetLeader(db *gorm.DB, userID, cartID string) error
}
