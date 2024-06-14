package data

import (
	"fmt"
	"log"

	"gorm.io/driver/clickhouse"
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
	log.Println("Migrating the database")
	return s.db.AutoMigrate(models...)
}

func (s *DataStore) Clean(models ...interface{}) error {
	log.Println("Dropping the tables")
	return s.db.Migrator().DropTable(models...)
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

func NewDataStore(cfg DbConfig) *DataStore {
	ds := InitializeDBFromConfig(cfg)
	err := ds.CheckConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return ds
}

func InitializeDBFromConfig(cfg DbConfig) *DataStore {
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
	case "clickhouse":
		dsn := fmt.Sprintf("tcp://%s:%s@%s:%d/%s",
			cfg.ClickHouse.Username, cfg.ClickHouse.Password, cfg.ClickHouse.Host, cfg.ClickHouse.Port, cfg.ClickHouse.DBName)
		db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported DB type: %s", cfg.Type)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database type: %s: %v", cfg.Type, err)
	}

	log.Printf("Database connection initialized to: %s\n", cfg.Type)

	return &DataStore{db: db}
}
