package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/synkube/app/evm-indexer/data"
	"github.com/synkube/app/evm-indexer/graphql/graph/model"
)

// Blocks is the resolver for the blocks field.
func (r *queryResolver) Blocks(ctx context.Context) ([]*model.Block, error) {
	blocks, err := r.BDS.GetAllBlocks()
	if err != nil {
		return nil, err
	}

	var result []*model.Block
	for _, block := range blocks {
		result = append(result, mapBlockToModel(block))
	}
	return result, nil
}

// Block is the resolver for the block field.
func (r *queryResolver) Block(ctx context.Context, id string) (*model.Block, error) {
	block, err := r.BDS.GetBlockByID(id)
	if err != nil {
		return nil, err
	}
	return mapBlockToModel(block), nil
}

// Transactions is the resolver for the transactions field.
func (r *queryResolver) Transactions(ctx context.Context) ([]*model.Transaction, error) {
	transactions, err := r.BDS.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	var result []*model.Transaction
	for _, tx := range transactions {
		result = append(result, mapTransactionToModel(tx))
	}
	return result, nil
}

// Transaction is the resolver for the transaction field.
func (r *queryResolver) Transaction(ctx context.Context, id string) (*model.Transaction, error) {
	tx, err := r.BDS.GetTransactionByID(id)
	if err != nil {
		return nil, err
	}
	return mapTransactionToModel(tx), nil
}

// Accounts is the resolver for the accounts field.
func (r *queryResolver) Accounts(ctx context.Context) ([]*model.Account, error) {
	accounts, err := r.BDS.GetAllAccounts()
	if err != nil {
		return nil, err
	}

	var result []*model.Account
	for _, account := range accounts {
		result = append(result, mapAccountToModel(account))
	}
	return result, nil
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context, address string) (*model.Account, error) {
	account, err := r.BDS.GetAccountByAddress(address)
	if err != nil {
		return nil, err
	}
	return mapAccountToModel(account), nil
}

// BlocksInRange is the resolver for the blocksInRange field.
func (r *queryResolver) BlocksInRange(ctx context.Context, startBlock string, endBlock string) ([]*model.Block, error) {
	panic(fmt.Errorf("not implemented: BlocksInRange - blocksInRange"))
}

// MissingBlocks is the resolver for the missingBlocks field.
func (r *queryResolver) MissingBlocks(ctx context.Context, startBlock string, endBlock string) ([]string, error) {
	panic(fmt.Errorf("not implemented: MissingBlocks - missingBlocks"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func mapBlockToModel(block *data.Block) *model.Block {
	return &model.Block{
		ID:              block.ID,
		Hash:            block.Hash,
		Number:          fmt.Sprint(block.Number),
		Timestamp:       block.Timestamp.String(),
		NumberOfTxs:     fmt.Sprint(block.NumberOfTxs),
		Miner:           block.Miner,
		ParentHash:      block.ParentHash,
		Difficulty:      block.Difficulty,
		TotalDifficulty: block.TotalDifficulty,
		Size:            fmt.Sprint(block.Size),
		GasUsed:         fmt.Sprint(block.GasUsed),
		GasLimit:        fmt.Sprint(block.GasLimit),
		Nonce:           block.Nonce,
		ExtraData:       block.ExtraData,
	}
}
func mapTransactionToModel(tx *data.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:               tx.ID,
		BlockHash:        tx.BlockHash,
		FromAddress:      tx.FromAddress,
		ToAddress:        fmt.Sprint(tx.ToAddress),
		Value:            tx.Value,
		Gas:              fmt.Sprint(tx.Gas),
		GasPrice:         tx.GasPrice,
		InputData:        tx.InputData,
		Nonce:            fmt.Sprint(tx.Nonce),
		TransactionIndex: fmt.Sprint(tx.TransactionIndex),
		Timestamp:        tx.Timestamp.String(),
	}
}
func mapAccountToModel(account *data.Account) *model.Account {
	return &model.Account{
		Address: account.Address,
		Balance: account.Balance,
	}
}
func mapTransactionsToModel(txs []*data.Transaction) []*model.Transaction {
	var result []*model.Transaction
	for _, tx := range txs {
		result = append(result, mapTransactionToModel(tx))
	}
	return result
}