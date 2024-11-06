package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"project-golang-crud/domains"
	"time"
)

type userRepository struct {
	db *gorm.DB
	genericRepo *GenericRepository
}

func NewUserRepository(db *gorm.DB) domains.UserRepository {
	genericRepo := &GenericRepository{db: db}
	return &userRepository{
		db: db,
		genericRepo: genericRepo,
	}
}

func (r *userRepository) CreateBookRequest(ctx context.Context,user *domains.User) error {
	return r.genericRepo.Create(ctx, user)
}

func (r *userRepository) GetAll(ctx context.Context) ([]domains.Book, error) {
	var books []domains.Book
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&books).Error // Mengambil hanya users yang deleted_at = null
	return books, err
}

func (r *userRepository) Create(ctx context.Context, user *domains.User) error {
	if user.DeletedAt != nil {
		return errors.New("user cannot be created because it is marked as deleted")
	}
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *userRepository) Update(ctx context.Context, user *domains.User) error {
	// Cek apakah user dengan deleted_at yang terisi sudah ada
	var existingUser domains.User
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NOT NULL", user.ID).First(&existingUser).Error; err == nil {
		return errors.New("cannot update user: user is marked as deleted")
	}
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	user := &domains.User{}

	if err := r.db.WithContext(ctx).First(user, "id = ?", id).Error; err != nil {
		return err
	}

	if user.DeletedAt != nil {
		return errors.New("user cannot be deleted because it is already marked as deleted")
	}

	now := time.Now()
	return r.db.WithContext(ctx).Model(&domains.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_at": now, "updated_at": gorm.Expr("updated_at")}).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domains.User, error) {
	var user domains.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domains.User, error) {
	var user domains.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err // Mengembalikan error jika tidak ditemukan
	}
	return &user, nil
}
