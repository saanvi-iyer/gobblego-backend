package table

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tableRepo struct{}

func NewMenuRepo() Repository {
	return &tableRepo{}
}

func (r *tableRepo) CreateTable(db *gorm.DB, table *models.Table) error {
	if table.TableID == uuid.Nil {
		table.TableID = uuid.New()
	}
	return db.Create(table).Error
}

func (r *tableRepo) GetTableByID(db *gorm.DB, id string) (*models.Table, error) {
	var table models.Table
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("menu_id = ?", uid).Where("deleted_at IS NULL").First(&table).Error
	return &table, err
}

func (r *tableRepo) ListTables(db *gorm.DB, limit, offset int) ([]models.Table, error) {
	var menus []models.Table
	err := db.Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&menus).Error
	return menus, err
}