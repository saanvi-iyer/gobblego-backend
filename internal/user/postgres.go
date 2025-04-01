package user

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct{}

func NewUserRepo() Repository {
	return &userRepo{}
}

func (r *userRepo) CreateUser(db *gorm.DB, user *models.User) error {
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}
	return db.Create(user).Error
}

func (r *userRepo) GetUserByID(db *gorm.DB, id string) (*models.User, error) {
	var user models.User
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("user_id = ?", uid).First(&user).Error
	return &user, err
}

func (r *userRepo) GetUsersByCartID(db *gorm.DB, cartID string) ([]models.User, error) {
	var users []models.User
	tid, err := uuid.Parse(cartID)
	if err != nil {
		return nil, err
	}
	err = db.Where("table_id = ?", tid).Find(&users).Error
	return users, err
}

func (r *userRepo) SetLeader(db *gorm.DB, userID, cartID string) error {
	tid, err := uuid.Parse(cartID)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&models.User{}).Where("table_id = ?", tid).Update("is_leader", false)
		return tx.Model(&models.User{}).Where("user_id = ?", uid).Update("is_leader", true).Error
	})
}