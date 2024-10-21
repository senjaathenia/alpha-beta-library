// repository/user_repository.go

package repository

import (
    "auth-user-api/models"

    "gorm.io/gorm"
)

type UserRepository interface {
    CreateUser(user *models.User) error
    GetUserByUsername(username string) (*models.User, error)
    GetUserByID(id string) (*models.User, error)
    UpdateUser(user *models.User) error
    DeleteUser(id string) error
    GetAllUsers() ([]*models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
    var users []*models.User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id string) error {
    return r.db.Model(&models.User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
