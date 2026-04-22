package database

import (
	"errors"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	Name           string
	WorldString    string
	OwnerPlayerID  uint
	WinnerPlayerID uint
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(name, worldString string, ownerPlayerID uint) (*Match, error) {
	match := &Match{
		Name:          name,
		WorldString:   worldString,
		OwnerPlayerID: ownerPlayerID,
	}
	result := r.db.Create(match)
	return match, result.Error
}

func (r *MatchRepository) FindByID(id uint) (*Match, error) {
	var match Match
	result := r.db.First(&match, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &match, result.Error
}

func (r *MatchRepository) FindAll() ([]Match, error) {
	var matches []Match
	result := r.db.Find(&matches)
	return matches, result.Error
}

func (r *MatchRepository) Update(id uint, updates map[string]interface{}) (*Match, error) {
	match, err := r.FindByID(id)
	if err != nil || match == nil {
		return nil, err
	}
	result := r.db.Model(match).Updates(updates)
	return match, result.Error
}

func (r *MatchRepository) Delete(id uint) error {
	result := r.db.Delete(&Match{}, id)
	return result.Error
}
