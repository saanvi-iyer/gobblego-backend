package menu

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type menuRepo struct{}

func NewMenuRepo() Repository {
	return &menuRepo{}
}

func (r *menuRepo) CreateMenu(db *gorm.DB, menu *models.Menu) error {
	if menu.ItemID == uuid.Nil {
		menu.ItemID = uuid.New()
	}
	return db.Create(menu).Error
}

func (r *menuRepo) GetMenuByID(db *gorm.DB, id string) (*models.Menu, error) {
	var menu models.Menu
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("menu_id = ?", uid).Where("deleted_at IS NULL").First(&menu).Error
	return &menu, err
}

func (r *menuRepo) ListMenus(db *gorm.DB, limit, offset int) ([]models.Menu, error) {
	var menus []models.Menu
	err := db.Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&menus).Error
	return menus, err
}

func (r *menuRepo) UpdateMenu(db *gorm.DB, menu *models.Menu) error {
	return db.Save(menu).Error
}
