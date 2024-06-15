package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"log"

	"github.com/synkube/app/blueprint/graphql/graph/model"
)

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.DataModel.GetUsers()
	if err != nil {
		return nil, err
	}

	var result []*model.User
	for _, user := range users {
		result = append(result, &model.User{
			ID:        fmt.Sprintf("%v", user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}
	return result, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	log.Println("User retrieval")
	user, err := r.DataModel.GetUserByID(fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        fmt.Sprintf("%v", id),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }