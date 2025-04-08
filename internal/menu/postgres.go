package menu

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
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
	err = db.Where("item_id = ?", uid).First(&menu).Error
	return &menu, err
}

func (r *menuRepo) ListMenus(db *gorm.DB, limit, offset int) ([]models.Menu, error) {
	var menus []models.Menu
	err := db.Limit(limit).Offset(offset).Find(&menus).Error
	return menus, err
}

func (r *menuRepo) UpdateMenu(db *gorm.DB, menu *models.Menu) error {
	return db.Save(menu).Error
}
