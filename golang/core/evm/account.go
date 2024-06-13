package evm

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// AccountCache is a thread-safe in-memory cache for account balances
type AccountCache struct {
	sync.RWMutex
	accounts map[common.Address]*big.Int
}

// NewAccountCache creates a new AccountCache
func NewAccountCache() *AccountCache {
	return &AccountCache{
		accounts: make(map[common.Address]*big.Int),
	}
}

// Get retrieves an account balance from the cache
func (ac *AccountCache) Get(addr common.Address) (*big.Int, bool) {
	ac.RLock()
	defer ac.RUnlock()
	balance, found := ac.accounts[addr]
	return balance, found
}

// Set adds an account balance to the cache
func (ac *AccountCache) Set(addr common.Address, balance *big.Int) {
	ac.Lock()
	defer ac.Unlock()
	ac.accounts[addr] = balance
}
