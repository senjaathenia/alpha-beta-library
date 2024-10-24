package repository

import (
	"project-golang-crud/domains"
	"time"
	"gorm.io/gorm"
	"errors"
)

type userRepository struct {
	db *gorm.DB 
}

func NewUserRepository(db *gorm.DB) domains.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]domains.User, error) {
	var users []domains.User
	err := r.db.Where("deleted_at IS NULL").Find(&users).Error // Mengambil hanya users yang deleted_at = null
	return users, err
}


func (r *userRepository) Create(user *domains.User) error {
	if user.DeletedAt != nil {
        return errors.New("user cannot be updated because it is marked as deleted")
    }
    return r.db.Create(&user).Error
}

func (r *userRepository) Update(user *domains.User) error {
	// Cek apakah user dengan deleted_at yang terisi sudah ada
	var existingUser domains.User
	if err := r.db.Where("id = ? AND deleted_at IS NOT NULL", user.ID).First(&existingUser).Error; err == nil {
		return errors.New("cannot update user: user is marked as deleted")
	}
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	user := &domains.User{}
    
    if err := r.db.First(user, "id = ?", id).Error; err != nil {
        return err 
    }
    
    if user.DeletedAt != nil {
        return errors.New("user cannot be deleted because it is already marked as deleted")
    }

    now := time.Now() 
	return r.db.Model(&domains.User{}).
    Where("id = ?", id).
    Updates(map[string]interface{}{"deleted_at": now, "updated_at": gorm.Expr("updated_at")}).Error
}

func (r *userRepository) GetByUsername(username string) (*domains.User, error) {
	var user domains.User
	if err :=r.db.Where("username = ?", username).First(&user).Error; err != nil{
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByID(id string) (*domains.User, error) {
	var user domains.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err // Mengembalikan error jika tidak ditemukan
	}
	return &user, nil
}

