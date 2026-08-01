package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/pkg/utils"

	"github.com/google/uuid"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByEmailOrUsername(email, username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("Email = ? OR Username = ?", email, username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("Email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := db.DB.Where("UserID = ?", utils.ToMSSQLUUIDString(id)).First(&user).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return db.DB.Model(&models.User{}).
		Where("UserID = ?", utils.ToMSSQLUUIDString(uuid.UUID(user.UserID))).
		Select("*").
		Omit("UserID", "CreatedAt").
		Updates(user).Error
}

