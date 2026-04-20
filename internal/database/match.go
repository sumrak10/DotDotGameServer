package database

import (
	"errors"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	Name             string
	WorldBuilderName string
	OwnerPlayerID    uint
	WinnerPlayerID   uint
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(name, worldBuilderName string, ownerPlayerID uint) (*Match, error) {
	user := &Match{
		Name:             name,
		WorldBuilderName: worldBuilderName,
		OwnerPlayerID:    ownerPlayerID,
	}
	result := r.db.Create(user)
	return user, result.Error
}

func (r *MatchRepository) FindByID(id uint) (*Match, error) {
	var user Match
	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *MatchRepository) FindAll() ([]Match, error) {
	var users []Match
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *MatchRepository) Update(id uint, updates map[string]interface{}) (*Match, error) {
	user, err := r.FindByID(id)
	if err != nil || user == nil {
		return nil, err
	}
	result := r.db.Model(user).Updates(updates)
	return user, result.Error
}

func (r *MatchRepository) Delete(id uint) error {
	result := r.db.Delete(&Match{}, id)
	return result.Error
}
