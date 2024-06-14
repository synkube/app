package data

import (
	"log"

	"github.com/synkube/app/blueprint/config"
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
	// Populate(ds)

	return ds
}

func Populate(ds *coreData.DataStore) {
	log.Println("Populating the database with sample data")
	ds.DB().FirstOrCreate(&User{Name: "John Doe", Email: "johndoe@example.com", Password: "password123"})
}
