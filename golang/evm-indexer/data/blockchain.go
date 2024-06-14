package data

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/synkube/app/core/data"
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
func (bds *BlockchainDataStore) SaveBlock(block *Block, transactions []*Transaction, accounts []*Account) error {
	log.Printf("Starting to save block number %d", block.Number)

	// Check if the block already exists
	exists, err := bds.blockExists(block.Number)
	if err != nil {
		log.Printf("Error checking existence of block number %d: %v", block.Number, err)
		return err
	} else if exists {
		log.Printf("Block number %d already exists, skipping save", block.Number)
		return nil
	}

	// Save transactions sequentially
	txsIds := make([]string, 0, len(transactions))
	for _, txn := range transactions {
		if err := bds.SaveTransaction(txn); err != nil {
			log.Printf("Error saving transaction %s: %v", txn.ID, err)
			return err
		}
		txsIds = append(txsIds, txn.ID)
	}

	// Save accounts sequentially
	for _, account := range accounts {
		if err := bds.SaveAccount(account); err != nil {
			log.Printf("Error saving account %s: %v", account.Address, err)
			return err
		}
	}

	// Update block details and save it
	block.NumberOfTxs = uint64(len(transactions))

	if err := bds.ds.DB().Save(block).Error; err != nil {
		log.Printf("Error saving block number %d: %v", block.Number, err)
		return err
	}

	log.Printf("Block number %d saved successfully", block.Number)
	return nil
}

// SaveTransaction saves a transaction to the database
func (bds *BlockchainDataStore) SaveTransaction(tx *Transaction) error {
	log.Printf("Starting to save transaction %s", tx.ID)
	if err := bds.ds.DB().Save(tx).Error; err != nil {
		log.Printf("Error saving transaction %s: %v", tx.ID, err)
		return err
	}
	log.Printf("Transaction %s saved successfully", tx.ID)
	return nil
}

// blockExists checks if a block with the given number exists in the database
func (bds *BlockchainDataStore) blockExists(number uint64) (bool, error) {
	var count int64
	if err := bds.ds.DB().Model(&Block{}).Where("number = ?", number).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveAccount saves an account to the database
func (bds *BlockchainDataStore) SaveAccount(account *Account) error {
	var count int64
	err := bds.ds.DB().Model(&Account{}).Where("address = ?", account.Address).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		// Account already exists, skip saving
		return nil
	}
	return bds.ds.DB().Save(account).Error
}

// GetLatestSavedBlock retrieves the latest saved block number from the database
func (bds *BlockchainDataStore) GetLatestSavedBlock() (uint64, error) {
	var count int64
	if err := bds.ds.DB().Model(&Block{}).Count(&count).Error; err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, nil // No blocks found
	}

	var block Block
	if err := bds.ds.DB().Order("number desc").First(&block).Error; err != nil {
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
	if len(missedBlocks) > 20 {
		log.Printf("Missing blocks: %v...", missedBlocks[:20])
	} else {
		log.Printf("Missing blocks: %v", missedBlocks)
	}

	return missedBlocks
}

// CreateBlockData creates a Block struct from the raw block data
func CreateBlockData(block *types.Block) *Block {
	return &Block{
		ID:     block.Hash().Hex(),
		Hash:   block.Hash().Hex(),
		Number: block.NumberU64(),
		// Transactions: 	block.Transactions(),
		NumberOfTxs:     uint64(block.Transactions().Len()),
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
	to := tx.To()
	if to == nil {
		to = &common.Address{}
	}

	return &Transaction{
		ID:          tx.Hash().Hex(),
		BlockHash:   block.Hash().Hex(),
		FromAddress: from.Hex(),
		ToAddress:   to.Hex(),
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
