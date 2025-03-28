package menu

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateMenu(db *gorm.DB, menu *models.Menu) error
	GetMenuByID(db *gorm.DB, id string) (*models.Menu, error)
	ListMenus(db *gorm.DB, limit, offset int) ([]models.Menu, error)
	UpdateMenu(db *gorm.DB, menu *models.Menu) error
}
