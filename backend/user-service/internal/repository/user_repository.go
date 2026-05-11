package repository

import (
	"github.com/deliverx/user-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindDrivers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("role = ? AND is_active = ?", model.RoleDriver, true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
