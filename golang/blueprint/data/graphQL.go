package data

import (
	"github.com/graphql-go/graphql"
)

var GraphQLSchema *graphql.Schema

func NewGraphQLSchema(dm *DataModel) *graphql.Schema {

	// Define the GraphQL object type for User
	var UserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Email": &graphql.Field{
				Type: graphql.String,
			},
			"CreatedAt": &graphql.Field{
				Type: graphql.String,
			},
			"UpdatedAt": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Define the user-related queries
	var UserQuery = graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(UserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dm.GetUsers(), nil
			},
		},
		"user": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if !ok {
					return nil, nil
				}
				return dm.GetUserByID(id), nil
			},
		},
	}

	// Combine the queries into the RootQuery
	var RootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": UserQuery["users"],
			"user":  UserQuery["user"],
			// Add other queries here
		},
	})

	// Define the schema
	GraphQLSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery,
	})
	if err != nil {
		panic(err)
	}
	return &GraphQLSchema
}
