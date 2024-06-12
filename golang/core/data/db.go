package data

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataStore struct {
	db *gorm.DB
}

func (s *DataStore) DB() *gorm.DB {
	return s.db
}

func (s *DataStore) Migrate(models ...interface{}) error {
	return s.db.AutoMigrate(models...)
}

func (s *DataStore) CheckConnection() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	} else {
		log.Println("Database connection is OK")
		return nil
	}
}

func InitializeDB(cfg DbConfig) *DataStore {
	var db *gorm.DB
	var err error
	switch cfg.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.SQLite.File), &gorm.Config{})
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported DB type: %s", cfg.Type)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database type %s: %v", cfg.Type, err)
	}

	log.Printf("Database connection initialized to %s\n", cfg.Type)

	return &DataStore{db: db}
}