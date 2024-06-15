package data

import (
	"log"
	"time"

	"github.com/synkube/app/blueprint/config"
	"github.com/synkube/app/core/common"
	coreData "github.com/synkube/app/core/data"
)

var models = []interface{}{
	&User{},
}

func Initialize(cfg *config.Config) *coreData.DataStore {
	ds := coreData.NewDataStore(cfg.DbConfig)
	if cfg.DbConfig.Clean {
		ds.Clean(models...)
	}
	ds.Migrate(models...)
	Populate(ds)

	return ds
}

func Populate(ds *coreData.DataStore) {
	log.Println("Populating the database with sample data")
	users := []User{
		{Name: "John Doe", Email: "johndoe@example.com", Password: "password123", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Jane Doe1", Email: "janedoe2@example.com", Password: "password1241153345", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Jane Doe2", Email: "janedoe1@example.com", Password: "password1tq23345"},
		{Name: "Jane Doe3", Email: "janedoe3@example.com", Password: "password1tew23345"},
	}
	common.PrettyPrint(users)
	for _, user := range users {
		ds.DB().Create(&user)
	}
}
