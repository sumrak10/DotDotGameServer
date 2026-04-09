package database

import "gorm.io/gorm"

type Match struct {
	gorm.Model
	Name string
}
