package data

import (
	"time"

	"github.com/synkube/app/core/data"
	"gorm.io/gorm"
)

type DataModel struct {
	ds *data.DataStore
	DB *gorm.DB
	UserModel
}

type UserModel interface {
	GetUsers() []User
	GetUserByID(id int) User
}

func NewDataModel(ds *data.DataStore) *DataModel {
	return &DataModel{
		ds: ds,
		DB: ds.DB(),
	}
}

func (dm *DataModel) GetUsers() []User {
	var users []User
	dm.DB.Find(&users)
	return users
}

func (dm *DataModel) GetUserByID(id int) User {
	var user User
	dm.DB.First(&user, id)
	return user
}

// User represents the User table in the database
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:100;not null"`
	Email     string         `gorm:"uniqueIndex;size:100;not null"`
	Password  string         `gorm:"size:100;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
