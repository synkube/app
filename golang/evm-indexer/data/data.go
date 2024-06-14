package data

import (
	"log"
	"time"

	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/evm-indexer/config"
)

// Block represents a block in the blockchain
type Block struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Hash      string    `json:"hash" gorm:"uniqueIndex"`
	Number    uint64    `json:"number"`
	Timestamp time.Time `json:"timestamp"`
	// Transactions    []Transaction `json:"transactions" gorm:"foreignKey:BlockHash;references:Hash"`
	// Transactions    []string `json:"transactions"` // Array of transaction IDs
	NumberOfTxs     uint64 `json:"numberOfTxs"`
	Miner           string `json:"miner"`
	ParentHash      string `json:"parentHash"`
	Difficulty      string `json:"difficulty"`
	TotalDifficulty string `json:"totalDifficulty"`
	Size            uint64 `json:"size"`
	GasUsed         uint64 `json:"gasUsed"`
	GasLimit        uint64 `json:"gasLimit"`
	Nonce           string `json:"nonce"`
	ExtraData       string `json:"extraData"`
}

// Transaction represents a transaction in the blockchain
type Transaction struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	BlockHash        string    `json:"blockHash" gorm:"index"`
	FromAddress      string    `json:"fromAddress"`
	ToAddress        string    `json:"toAddress"`
	Value            string    `json:"value"`
	Gas              uint64    `json:"gas"`
	GasPrice         string    `json:"gasPrice"`
	InputData        string    `json:"inputData"`
	Nonce            uint64    `json:"nonce"`
	TransactionIndex uint64    `json:"transactionIndex"`
	Timestamp        time.Time `json:"timestamp"`
}

// Account represents an account in the blockchain
type Account struct {
	Address string `json:"address" gorm:"primaryKey"`
	Balance string `json:"balance"`
	// Transactions []Transaction `json:"transactions" gorm:"foreignKey:FromAddress;references:Address"`
}

var models = []interface{}{
	&Account{},
	&Block{},
	&Transaction{},
}

func Initialize(cfg *config.Config) *coreData.DataStore {
	var ds *coreData.DataStore
	if cfg.DbConfig.Type != "" {
		ds = coreData.InitializeDBConn(cfg.DbConfig)
		ds.CheckConnection()
		if cfg.Indexer.Clean {
			ds.Clean(models...)
		}
		ds.Migrate(models...)
		// Populate(ds)
	} else {
		log.Println("No database configuration found")
	}
	return ds
}

func Populate(ds *coreData.DataStore) {
	log.Println("Populating the database with sample data")
}
