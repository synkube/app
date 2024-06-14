package indexer

import (
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/core/evm"
	"github.com/synkube/app/evm-indexer/config"
	"github.com/synkube/app/evm-indexer/data"
)

// Worker function for goroutines to index blocks
func worker(id int, bm *BlockManager, rpcClient *RPCClient, bds *data.BlockchainDataStore, cache *evm.AccountCache, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Starting", id)
	for {
		blockNumber, ok := bm.GetNextBlock()
		if !ok {
			log.Printf("Worker %d: No more blocks to process", id)
			return
		}
		log.Printf("Worker %d: Indexing block %d", id, blockNumber)
		err := indexBlock(rpcClient, blockNumber, bds, cache)
		if err != nil {
			log.Printf("Worker %d: Error indexing block %d: %v", id, blockNumber, err)
			// bm.AddMissedBlock(blockNumber) // TODO - Add back missed block
		} else {
			log.Printf("Worker %d: Successfully indexed block %d", id, blockNumber)
		}
	}
}

// indexBlock retrieves and processes a block
func indexBlock(rpcClient *RPCClient, blockNumber int, bds *data.BlockchainDataStore, cache *evm.AccountCache) error {
	log.Printf("Indexing block %d", blockNumber)
	block, err := rpcClient.GetBlockWithRetry(uint64(blockNumber))
	// block, err := evm.GetBlock(rpcClient, big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Printf("Error retrieving block %d: %v", blockNumber, err)
		return err
	}

	transactions, err := processTransactions(block, bds, cache, rpcClient)
	if err != nil {
		log.Printf("Error processing transactions for block %d: %v", blockNumber, err)
		return err
	}

	// blockTD, err := rpcClient.BlockByHash(context.Background(), block.Hash())
	blockTD, err := rpcClient.GetBlockByHashWithRetry(block.Hash())
	if err != nil {
		log.Printf("Failed to retrieve total difficulty for block %d: %v", blockNumber, err)
		return fmt.Errorf("failed to retrieve total difficulty: %v", err)
	}

	blockData := data.CreateBlockData(blockTD)

	err = bds.SaveBlock(blockData, transactions)
	if err != nil {
		log.Printf("Failed to save block %d: %v", blockNumber, err)
		return fmt.Errorf("failed to save block %d: %v", blockNumber, err)
	}

	log.Printf("Successfully indexed block %d", blockNumber)
	return nil
}

// processTransactions processes transactions within a block
func processTransactions(block *types.Block, bds *data.BlockchainDataStore, cache *evm.AccountCache, rpcClient *RPCClient) ([]*data.Transaction, error) {
	log.Printf("Processing transactions for block %d", block.Number().Uint64())
	var transactions []*data.Transaction

	for _, tx := range evm.GetTransactions(block) {
		txData, err := data.CreateTransactionData(tx, block)
		if err != nil {
			log.Printf("Failed to create transaction %s: %v", tx.Hash().Hex(), err)
			return nil, fmt.Errorf("failed to create transaction %s: %v", tx.Hash().Hex(), err)
		}

		err = bds.SaveTransaction(txData)
		if err != nil {
			log.Printf("Failed to save transaction %s: %v", tx.Hash().Hex(), err)
			return nil, fmt.Errorf("failed to save transaction %s: %v", tx.Hash().Hex(), err)
		}

		log.Printf("Successfully processed transaction %s", tx.Hash().Hex())
		transactions = append(transactions, txData)

		err = processAccounts(tx, cache, bds, rpcClient)
		if err != nil {
			log.Printf("Error processing accounts for transaction %s: %v", tx.Hash().Hex(), err)
			return nil, err
		}
	}

	log.Printf("Successfully processed transactions for block %d", block.Number().Uint64())
	return transactions, nil
}

// processAccounts processes accounts involved in a transaction
func processAccounts(tx *types.Transaction, cache *evm.AccountCache, bds *data.BlockchainDataStore, rpcClient *RPCClient) error {
	log.Printf("Processing accounts for transaction %s", tx.Hash().Hex())
	chainID := tx.ChainId()
	signer := types.NewEIP155Signer(chainID)

	// Extract the sender address from the transaction
	from, err := types.Sender(signer, tx)
	if err != nil {
		log.Printf("Failed to extract sender for transaction %s: %v", tx.Hash().Hex(), err)
		return fmt.Errorf("failed to extract sender: %v", err)
	}

	accounts := []common.Address{from}
	if tx.To() != nil {
		accounts = append(accounts, *tx.To())
	}

	for _, addr := range accounts {
		log.Printf("Processing account %s", addr.Hex())
		_, found := cache.Get(addr)
		if !found {
			balance, err := rpcClient.GetBalanceWithRetry(addr)
			if err != nil {
				log.Printf("Failed to retrieve account balance for %s: %v", addr.Hex(), err)
				return fmt.Errorf("failed to retrieve account balance for %s: %v", addr.Hex(), err)
			}

			accountData := data.CreateAccountData(addr, balance)

			err = bds.SaveAccount(accountData)
			if err != nil {
				log.Printf("Failed to save account %s: %v", addr.Hex(), err)
				return fmt.Errorf("failed to save account %s: %v", addr.Hex(), err)
			}

			cache.Set(addr, balance)
			log.Printf("Successfully processed account %s", addr.Hex())
		}
	}

	return nil
}

// StartIndexing initializes the process
func StartIndexing(chainConfig coreData.Chain, ds *coreData.DataStore, indexerConfig config.Indexer) error {
	log.Println("Starting indexing process...")
	// rpcManager := NewRPCManager(chainConfig)
	rpcClient, err := NewRPCClient(chainConfig.RPCs, indexerConfig.MaxRetries)
	if err != nil {
		log.Printf("Failed to create RPC client: %v", err)
		return fmt.Errorf("failed to create RPC client: %v", err)
	}

	// Initialize the BlockchainDataStore
	bds := data.NewBlockchainDataStore(ds)

	// Get the latest saved block from the data store
	latestSavedBlock, err := bds.GetLatestSavedBlock()
	if err != nil {
		log.Printf("Failed to get latest saved block: %v", err)
		return fmt.Errorf("failed to get latest saved block: %v", err)
	}

	// Ensure that the latest saved block is within the range of start and end blocks
	if latestSavedBlock <= uint64(indexerConfig.StartBlock) {
		latestSavedBlock = uint64(indexerConfig.StartBlock)
	}
	log.Printf("Starting from block %d", latestSavedBlock)

	// Identify missing blocks from startBlock to latestSavedBlock
	missedBlocks := bds.IdentifyMissingBlocks(uint64(indexerConfig.StartBlock), latestSavedBlock)

	// Create BlockManager with missing blocks from start to latest saved block
	blockManager := NewBlockManager(int(latestSavedBlock), indexerConfig.EndBlock)
	blockManager.AddMissedBlocks(missedBlocks)

	accountCache := evm.NewAccountCache()

	// Distribute the load across multiple goroutines
	var wg sync.WaitGroup
	numWorkers := indexerConfig.MaxWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, blockManager, rpcClient, bds, accountCache, &wg)
	}
	wg.Wait()
	log.Println("Indexing process completed")
	return nil
}
