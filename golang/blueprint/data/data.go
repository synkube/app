package data

import (
	"log"
	"time"

	"github.com/synkube/app/blueprint/config"
	coreData "github.com/synkube/app/core/data"
	"gorm.io/gorm"
)

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

func Initialize(cfg *config.Config) *coreData.DataStore {
	var ds *coreData.DataStore
	if cfg.DbConfig.Type != "" {
		ds = coreData.InitializeDB(cfg.DbConfig)
		ds.CheckConnection()
		Populate(ds)
	} else {
		log.Println("No database configuration found")
	}
	return ds
}

func Populate(ds *coreData.DataStore) {
	ds.Migrate(&User{})
	log.Println("Migrating the database")

	ds.DB().FirstOrCreate(&User{Name: "John Doe", Email: "johndoe@example.com", Password: "password123"})
	log.Println("Populating the database with sample data")
}
