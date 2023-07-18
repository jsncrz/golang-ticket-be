package ticket

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repo Repository) ticketCollection() *mongo.Collection {
	return repo.client.Database("APIintegration").Collection("tickets")
}

func (repo Repository) CreateTicket(ctx context.Context, ticket Ticket) error {
	_, err := repo.ticketCollection().InsertOne(ctx, ticket)
	return err
}

func updateTicketBson(t Ticket) bson.D {
	t.AddCompletedAtIfCompleted()
	return bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: t.Title},
		{Key: "desciption", Value: t.Description},
		{Key: "position", Value: t.Position},
		{Key: "updated_at", Value: time.Now()},
		{Key: "status", Value: t.Status},
		{Key: "priority", Value: t.Priority},
		{Key: "completed_at", Value: t.CompletedAt},
	},
	}}
}

func (repo Repository) UpdateTicket(ctx context.Context, ticket Ticket) error {
	_, err := repo.ticketCollection().UpdateByID(ctx, ticket.ID, updateTicketBson(ticket))
	return err
}

func (repo Repository) GetTickets(ctx context.Context) (Tickets, error) {
	var tickets Tickets
	opts := options.Find().SetSort(sortByPosition())
	results, err := repo.ticketCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ticket Ticket
		if err = results.Decode(&ticket); err != nil {
			return nil, err
		}

		tickets = append(tickets, &ticket)
	}
	return tickets, nil
}

func groupByStatus() bson.D {
	return bson.D{{
		Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{
				{Key: "$sum", Value: 1}},
			},
		},
	}}
}

func filterByStatus(status Status) bson.D {
	return bson.D{
		{Key: "status", Value: bson.D{{
			Key: "$eq", Value: status,
		}},
		},
	}
}
func sortByPosition() bson.D {
	return bson.D{
		{Key: "position", Value: 1},
	}
}

func (repo Repository) GetLastTicketByStatus(ctx context.Context, status Status) (*Ticket, error) {
	filter := filterByStatus(status)
	opts := options.Find().SetSort(sortByPosition()).SetLimit(2)
	results, err := repo.ticketCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer results.Close(ctx)
	var tickets []Ticket
	err2 := results.All(ctx, &tickets)
	if len(tickets) > 0 {
		return &tickets[0], nil
	} else {
		return nil, err2
	}
}
