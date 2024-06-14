package indexer

import (
	"fmt"
	"log"
	"sync"

	goEthCommon "github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"
	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/core/evm"
	"github.com/synkube/app/evm-indexer/config"
	"github.com/synkube/app/evm-indexer/data"
)

// Worker function for goroutines to index blocks
func worker(id int, bm *BlockManager, rpcClient *RPCClient, bds *data.BlockchainDataStore, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d: Starting", id)
	for {
		blockNumber, ok := bm.GetNextBlock()
		if !ok {
			log.Printf("Worker %d: No more blocks to process", id)
			return
		}
		log.Printf("Worker %d: Indexing block %d", id, blockNumber)
		err := indexBlock(rpcClient, blockNumber, bds)
		if err != nil {
			log.Printf("Worker %d: Error indexing block %d: %v", id, blockNumber, err)
			// bm.AddMissedBlock(blockNumber) // TODO - Add back missed block
		} else {
			log.Printf("Worker %d: Successfully indexed block %d", id, blockNumber)
		}
	}
}

// indexBlock retrieves and processes a block
func indexBlock(rpcClient *RPCClient, blockNumber int, bds *data.BlockchainDataStore) error {
	log.Printf("Indexing block %d", blockNumber)
	block, err := rpcClient.GetBlockWithRetry(uint64(blockNumber))
	if err != nil {
		log.Printf("Error retrieving block %d: %v", blockNumber, err)
		return err
	}

	goEthTxs := evm.GetTransactions(block)
	transactions, err := processTransactions(block, goEthTxs)
	if err != nil {
		log.Printf("Error processing transactions for block %d: %v", blockNumber, err)
		return err
	}

	accounts, err := processAccounts(goEthTxs)
	if err != nil {
		log.Printf("Error processing accounts for block %d: %v", blockNumber, err)
		return err
	}

	accountsWithBalance, err := retrieveAccountsWithBalance(rpcClient, accounts)
	if err != nil {
		log.Printf("Error retrieving accounts with balance for block %d: %v", blockNumber, err)
		return err
	}

	blockData := data.CreateBlockData(block)

	err = bds.SaveBlock(blockData, transactions, accountsWithBalance)
	if err != nil {
		log.Printf("Failed to save block %d: %v", blockNumber, err)
		return fmt.Errorf("failed to save block %d: %v", blockNumber, err)
	}

	log.Printf("Successfully indexed block %d", blockNumber)
	return nil
}

// processTransactions processes transactions within a block
func processTransactions(block *goEthTypes.Block, txs []*goEthTypes.Transaction) ([]*data.Transaction, error) {
	log.Printf("Processing transactions for block %d", block.Number().Uint64())
	var transactions []*data.Transaction = make([]*data.Transaction, 0)

	for _, tx := range txs {
		txData, err := data.CreateTransactionData(tx, block)
		if err != nil {
			log.Printf("Failed to create transaction %s: %v", tx.Hash().Hex(), err)
			return nil, fmt.Errorf("failed to create transaction %s: %v", tx.Hash().Hex(), err)
		}
		transactions = append(transactions, txData)
	}
	return transactions, nil
}

// processAccounts processes accounts involved in a transaction
func processAccounts(txs []*goEthTypes.Transaction) ([]goEthCommon.Address, error) {
	accounts := make(map[string]goEthCommon.Address)
	for _, tx := range txs {
		// Assuming From() method returns the sender's address (common.Address)
		from, err := goEthTypes.Sender(goEthTypes.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return nil, fmt.Errorf("failed to get sender address: %w", err)
		}
		accounts[from.Hex()] = from

		// Get the recipient's address
		to := tx.To()
		if to != nil {
			accounts[to.Hex()] = *to
		}
	}

	uniqueAddresses := make([]goEthCommon.Address, 0, len(accounts))
	for _, addr := range accounts {
		uniqueAddresses = append(uniqueAddresses, addr)
	}

	return uniqueAddresses, nil
}

func retrieveAccountsWithBalance(rpcClient *RPCClient, accounts []goEthCommon.Address) ([]*data.Account, error) {
	accountsWithBalance := make([]*data.Account, 0, len(accounts))

	for _, account := range accounts {
		balance, err := rpcClient.GetBalanceWithRetry(account)
		if err != nil {
			log.Printf("Failed to retrieve account balance for %s: %v", account.Hex(), err)
			return nil, fmt.Errorf("failed to retrieve account balance for %s: %v", account.Hex(), err)
		}
		accountData := data.CreateAccountData(account, balance)

		accountsWithBalance = append(accountsWithBalance, accountData)
	}

	return accountsWithBalance, nil
}

// StartIndexing initializes the process
func StartIndexing(chainConfig coreData.Chain, ds *coreData.DataStore, indexerConfig config.Indexer) error {
	log.Println("## Starting indexing process...")
	log.Println("Setup RPC client")
	rpcClient, err := NewRPCClient(chainConfig.RPCs, indexerConfig.MaxRetries)
	if err != nil {
		log.Printf("Failed to create RPC client: %v", err)
		return fmt.Errorf("failed to create RPC client: %v", err)
	}

	// Initialize the BlockchainDataStore
	log.Println("Initialize BlockchainDataStore")
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
	log.Println("Identify missing blocks")
	missedBlocks := bds.IdentifyMissingBlocks(uint64(indexerConfig.StartBlock), latestSavedBlock)

	// Create BlockManager with missing blocks from start to latest saved block
	blockManager := NewBlockManager(int(latestSavedBlock), indexerConfig.EndBlock)
	blockManager.AddMissedBlocks(missedBlocks)

	// Distribute the load across multiple goroutines
	var wg sync.WaitGroup
	numWorkers := indexerConfig.MaxWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, blockManager, rpcClient, bds, &wg)
	}
	wg.Wait()
	log.Println("Indexing process completed")
	return nil
}
