package database

import (
	"errors"

	"gorm.io/gorm"
)

type UserView struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
}

type UserSensitiveView struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;column:username"`
	Email    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
}

func (u *User) ToView() *UserView {
	return &UserView{
		ID:       u.ID,
		UserName: u.Username,
	}
}

func (u *User) ToSensitiveView() *UserSensitiveView {
	return &UserSensitiveView{
		ID:       u.ID,
		UserName: u.Username,
		Email:    u.Email,
	}
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username, email, password string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}
	result := r.db.Create(user)
	return user, result.Error
}

func (r *UserRepository) FindByID(id uint) (*User, error) {
	var user User
	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *UserRepository) FindByUserName(username string) (*User, error) {
	var user User
	result := r.db.First(&user, "username = ?", username)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := r.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *UserRepository) FindAll() ([]User, error) {
	var users []User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepository) Update(id uint, updates map[string]interface{}) (*User, error) {
	user, err := r.FindByID(id)
	if err != nil || user == nil {
		return nil, err
	}
	result := r.db.Model(user).Updates(updates)
	return user, result.Error
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&User{}, id)
	return result.Error
}
