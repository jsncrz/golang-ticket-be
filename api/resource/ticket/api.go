package ticket

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	repo *Repository
}

func New(c *mongo.Client) *API {
	return &API{
		repo: NewRepository(c),
	}
}
