package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlock retrieves a block by its number
func GetBlock(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}
	return block, nil
}

// GetTransactions retrieves transactions from a specific block
func GetTransactions(block *types.Block) []*types.Transaction {
	return block.Transactions()
}

// GetAccount retrieves the account balance
func GetAccount(client *ethclient.Client, accountAddress common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account balance: %v", err)
	}
	return balance, nil
}
