package indexer

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ava-labs/coreth/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	coreData "github.com/synkube/app/core/data"
)

type RPCManager struct {
	primary   coreData.RPC
	auxiliary []coreData.RPC
	mu        sync.Mutex
	current   int
}

// NewRPCManager creates a new RPCManager
func NewRPCManager(chainConfig coreData.Chain) *RPCManager {
	var primary coreData.RPC
	var auxiliary []coreData.RPC
	for _, rpc := range chainConfig.RPCs {
		if rpc.Type == "primary" {
			primary = coreData.RPC{URL: rpc.URL, Type: rpc.Type}
		} else {
			auxiliary = append(auxiliary, coreData.RPC{URL: rpc.URL, Type: rpc.Type})
		}
	}
	return &RPCManager{
		primary:   primary,
		auxiliary: auxiliary,
	}
}

// GetRPC returns the current RPC to use
func (rm *RPCManager) GetRPC() string {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.current == 0 {
		return rm.primary.URL
	}
	return rm.auxiliary[rm.current-1].URL
}

// Handle429 switches to the next available RPC on a 429 error
func (rm *RPCManager) Handle429() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.current++
	if rm.current >= len(rm.auxiliary)+1 {
		rm.current = 0
	}
}

type BlockManager struct {
	sync.Mutex
	currentBlock int
	maxBlock     int
	missedBlocks map[int]struct{}
}

// NewBlockManager creates a new BlockManager
func NewBlockManager(startBlock, maxBlock int) *BlockManager {
	return &BlockManager{
		currentBlock: startBlock,
		maxBlock:     maxBlock,
		missedBlocks: make(map[int]struct{}),
	}
}

// GetNextBlock returns the next block to be indexed
func (bm *BlockManager) GetNextBlock() (int, bool) {
	bm.Lock()
	defer bm.Unlock()

	if len(bm.missedBlocks) > 0 {
		for block := range bm.missedBlocks {
			delete(bm.missedBlocks, block)
			return block, true
		}
	}

	if bm.currentBlock <= bm.maxBlock {
		block := bm.currentBlock
		bm.currentBlock++
		return block, true
	}

	return 0, false
}

// AddMissedBlock adds a missed block to be re-indexed
func (bm *BlockManager) AddMissedBlock(block int) {
	bm.Lock()
	defer bm.Unlock()
	bm.missedBlocks[block] = struct{}{}
}

// AddMissedBlock adds a missed block to be re-indexed
func (bm *BlockManager) AddMissedBlocks(blocks []int) {
	bm.Lock()
	defer bm.Unlock()
	for _, block := range blocks {
		bm.missedBlocks[block] = struct{}{}
	}
}

type RPCClient struct {
	rpcs          []coreData.RPC
	currentRPCIdx int
	maxRetries    int
	client        *ethclient.Client
	mutex         sync.Mutex // To ensure thread-safe access to currentRPCIdx
}

// NewRPCClient creates a new RPCClient instance
func NewRPCClient(rpcs []coreData.RPC, maxRetries int) (*RPCClient, error) {
	rpcClient := &RPCClient{
		rpcs:       rpcs,
		maxRetries: maxRetries,
	}

	client, err := rpcClient.connectWithRetry()
	if err != nil {
		return nil, err
	}
	rpcClient.client = client

	return rpcClient, nil
}

func (rpcClient *RPCClient) connectWithRetry() (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error
	for i := 0; i < len(rpcClient.rpcs); i++ {
		url := rpcClient.rpcs[i].URL
		for retry := 0; retry < rpcClient.maxRetries; retry++ {
			client, err = ethclient.Dial(url)
			if err == nil {
				// Successfully connected
				rpcClient.mutex.Lock()
				rpcClient.currentRPCIdx = i
				rpcClient.mutex.Unlock()
				return client, nil
			}
			// Check if the error is due to rate limiting (HTTP 429)
			if rpcErr, ok := err.(*rpc.HTTPError); ok && rpcErr.StatusCode == 429 {
				log.Printf("RPC %s rate limited (429). Switching to next RPC URL.", url)
				break
			}
			log.Printf("Failed to connect to Ethereum client (%s): %v. Retrying (%d/%d)...", url, err, retry+1, rpcClient.maxRetries)
			time.Sleep(2 * time.Second)
		}
	}
	return nil, fmt.Errorf("failed to connect to any Ethereum client after %d retries", rpcClient.maxRetries*len(rpcClient.rpcs))
}

func (rpcClient *RPCClient) switchRPC() error {
	rpcClient.mutex.Lock()
	defer rpcClient.mutex.Unlock()

	rpcClient.currentRPCIdx++
	if rpcClient.currentRPCIdx >= len(rpcClient.rpcs) {
		rpcClient.currentRPCIdx = 0
	}
	client, err := rpcClient.connectWithRetry()
	if err != nil {
		return err
	}
	rpcClient.client = client
	return nil
}

func (rpcClient *RPCClient) retry(f func() error) error {
	var err error
	for retry := 0; retry < rpcClient.maxRetries; retry++ {
		err = f()
		if err == nil {
			// Successfully executed the function
			return nil
		}
		// Check if the error is due to rate limiting (HTTP 429)
		if rpcErr, ok := err.(*rpc.HTTPError); ok && rpcErr.StatusCode == 429 {
			log.Printf("RPC rate limited (429). Switching to next RPC URL.")
			err = rpcClient.switchRPC()
			if err != nil {
				return err
			}
			continue
		}
		log.Printf("Error executing function: %v. Retrying...", err)
		time.Sleep(2 * time.Second)
	}
	return err
}

func (rpcClient *RPCClient) GetBlockWithRetry(blockNumber uint64) (*types.Block, error) {
	var block *types.Block
	err := rpcClient.retry(func() error {
		var err error
		block, err = rpcClient.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		return err
	})
	return block, err
}

func (rpcClient *RPCClient) GetBlockByHashWithRetry(hash common.Hash) (*types.Block, error) {
	var block *types.Block
	err := rpcClient.retry(func() error {
		var err error
		block, err = rpcClient.client.BlockByHash(context.Background(), hash)
		return err
	})
	return block, err
}

func (rpcClient *RPCClient) GetBalanceWithRetry(account common.Address) (*big.Int, error) {
	var balance *big.Int
	err := rpcClient.retry(func() error {
		var err error
		balance, err = rpcClient.client.BalanceAt(context.Background(), account, nil)
		return err
	})
	return balance, err
}
