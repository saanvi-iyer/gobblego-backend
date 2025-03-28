package user

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(db *gorm.DB, user *models.User) error
	GetUserByID(db *gorm.DB, id string) (*models.User, error)
	GetUsersByTableID(db *gorm.DB, tableID string) ([]models.User, error)
	SetLeader(db *gorm.DB, userID, tableID string) error
}
