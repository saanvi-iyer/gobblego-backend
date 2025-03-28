package table

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateTable(db *gorm.DB, menu *models.Table) error
	GetTableByID(db *gorm.DB, id string) (*models.Table, error)
	ListTables(db *gorm.DB, limit, offset int) ([]models.Table, error)
}
