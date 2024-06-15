package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import "github.com/synkube/app/evm-indexer/data"

type Resolver struct {
	BDS *data.BlockchainDataStore
}
