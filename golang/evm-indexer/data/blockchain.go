package data

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/synkube/app/core/data"
	"gorm.io/gorm"
)

// BlockchainDataStore wraps the GORM DB instance to provide higher-level operations
type BlockchainDataStore struct {
	ds *data.DataStore
}

// NewBlockchainDataStore creates a new BlockchainDataStore
func NewBlockchainDataStore(ds *data.DataStore) *BlockchainDataStore {
	return &BlockchainDataStore{ds: ds}
}

// SaveBlock saves a block and its transactions to the database
func (bds *BlockchainDataStore) SaveBlock(block *Block, transactions []*Transaction) error {
	err := bds.ds.DB().Transaction(func(db *gorm.DB) error {
		var existingBlock Block
		err := db.Where("number = ?", block.Number).First(&existingBlock).Error
		if err == nil {
			// Block already exists, skip saving
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		for _, txn := range transactions {
			if err := db.Save(txn).Error; err != nil {
				return err
			}
		}
		if err := db.Save(block).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// SaveTransaction saves a transaction to the database
func (bds *BlockchainDataStore) SaveTransaction(tx *Transaction) error {
	var existingTx Transaction
	err := bds.ds.DB().Where("id = ?", tx.ID).First(&existingTx).Error
	if err == nil {
		// Transaction already exists, skip saving
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return bds.ds.DB().Save(tx).Error
}

// SaveAccount saves an account to the database
func (bds *BlockchainDataStore) SaveAccount(account *Account) error {
	var existingAccount Account
	err := bds.ds.DB().Where("address = ?", account.Address).First(&existingAccount).Error
	if err == nil {
		// Transaction already exists, skip saving
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return bds.ds.DB().Save(account).Error
}

// GetLatestSavedBlock retrieves the latest saved block number from the database
func (bds *BlockchainDataStore) GetLatestSavedBlock() (uint64, error) {
	var block Block
	if err := bds.ds.DB().Order("number desc").First(&block).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // No blocks found, return 0
		}
		return 0, err
	}
	return block.Number, nil
}

// GetAllBlockNumbers retrieves all block numbers from the database
func (bds *BlockchainDataStore) GetAllBlockNumbers() ([]uint64, error) {
	var blockNumbers []uint64
	rows, err := bds.ds.DB().Model(&Block{}).Select("number").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var number uint64
		if err := rows.Scan(&number); err != nil {
			return nil, err
		}
		blockNumbers = append(blockNumbers, number)
	}

	return blockNumbers, nil
}

// GetBlockNumbersInRange retrieves all block numbers in the specified range from the database
func (bds *BlockchainDataStore) GetBlockNumbersInRange(startBlock, endBlock uint64) ([]uint64, error) {
	var blockNumbers []uint64
	rows, err := bds.ds.DB().Model(&Block{}).Select("number").Where("number >= ? AND number <= ?", startBlock, endBlock).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var number uint64
		if err := rows.Scan(&number); err != nil {
			return nil, err
		}
		blockNumbers = append(blockNumbers, number)
	}

	return blockNumbers, nil
}

func (bds *BlockchainDataStore) BlockExists(blockNumber uint64) (bool, error) {
	var block Block
	err := bds.ds.DB().Where("number = ?", blockNumber).First(&block).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// identifyMissingBlocks identifies any missing blocks between startBlock and latestSavedBlock
func (bds *BlockchainDataStore) IdentifyMissingBlocks(startBlock, latestSavedBlock uint64) []int {
	log.Printf("Identifying missing blocks between %d and %d", startBlock, latestSavedBlock)
	allBlockNumbers, err := bds.GetBlockNumbersInRange(startBlock, latestSavedBlock)
	if err != nil {
		log.Printf("Error retrieving block numbers in range: %v", err)
		return nil
	}

	blockNumberSet := make(map[uint64]struct{}, len(allBlockNumbers))
	for _, number := range allBlockNumbers {
		blockNumberSet[number] = struct{}{}
	}

	missedBlocks := []int{}
	for blockNumber := startBlock; blockNumber <= latestSavedBlock; blockNumber++ {
		if _, exists := blockNumberSet[blockNumber]; !exists {
			missedBlocks = append(missedBlocks, int(blockNumber))
		}
	}
	log.Printf("Found %d missing blocks", len(missedBlocks))

	return missedBlocks
}

// CreateBlockData creates a Block struct from the raw block data
func CreateBlockData(block *types.Block) *Block {
	return &Block{
		ID:              block.Hash().Hex(),
		Hash:            block.Hash().Hex(),
		Number:          block.NumberU64(),
		Timestamp:       time.Unix(int64(block.Time()), 0),
		Miner:           block.Coinbase().Hex(),
		ParentHash:      block.ParentHash().Hex(),
		Difficulty:      block.Difficulty().String(),
		TotalDifficulty: block.Difficulty().String(),
		Size:            block.Size(),
		GasUsed:         block.GasUsed(),
		GasLimit:        block.GasLimit(),
		Nonce:           fmt.Sprintf("%d", block.Nonce()),
		ExtraData:       fmt.Sprintf("%x", block.Extra()),
	}
}

// CreateTransactionData creates a Transaction struct from the raw transaction data
func CreateTransactionData(tx *types.Transaction, block *types.Block) (*Transaction, error) {
	chainID := tx.ChainId()
	signer := types.NewEIP155Signer(chainID)

	// Extract the sender address from the transaction
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sender: %v", err)
	}

	return &Transaction{
		ID:          tx.Hash().Hex(),
		BlockHash:   block.Hash().Hex(),
		FromAddress: from.Hex(),
		ToAddress:   tx.To().Hex(),
		Value:       tx.Value().String(),
		Gas:         tx.Gas(),
		GasPrice:    tx.GasPrice().String(),
		InputData:   fmt.Sprintf("%x", tx.Data()),
		Nonce:       tx.Nonce(),
		Timestamp:   time.Unix(int64(block.Time()), 0),
	}, nil
}

// CreateAccountData creates an Account struct from the raw account data
func CreateAccountData(address common.Address, balance *big.Int) *Account {
	return &Account{
		Address: address.Hex(),
		Balance: balance.String(),
	}
}
