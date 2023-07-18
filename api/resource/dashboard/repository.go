package dashboard

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func groupByYearMonth() bson.D {
	return bson.D{{
		Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "month", Value: bson.D{{Key: "$month", Value: "$start_time"}}},
				{Key: "year", Value: bson.D{{Key: "$year", Value: "$start_time"}}},
			}},
			{Key: "numberOfTickets", Value: bson.D{
				{Key: "$sum", Value: 1}},
			},
		},
	}}
}
func (repo Repository) GetMonthlyTicketsCount(ctx context.Context) (MonthlyTickets, error) {
	var monthlyTickets MonthlyTickets
	groupStage := groupByYearMonth()

	results, err := repo.ticketCollection().Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return nil, err
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var monthlyTicket MonthlyTicket
		if err = results.Decode(&monthlyTicket); err != nil {
			return nil, err
		}

		monthlyTickets = append(monthlyTickets, &monthlyTicket)
	}
	return monthlyTickets, nil
}
