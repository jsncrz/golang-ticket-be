package ticket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	ticket := &Ticket{}

	if err := json.NewDecoder(r.Body).Decode(ticket); err != nil {
		log.Print(err)
		return
	}

	newTicket := Ticket{
		ID:          primitive.NewObjectID(),
		Title:       ticket.Title,
		Description: ticket.Description,
		Position:    ticket.Position,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	newTicket.AddCompletedAtIfCompleted()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	lastTicket, err := a.repo.GetLastTicketByStatus(ctx, ticket.Status)
	if err != nil {
		log.Print(err)
	}
	log.Print(lastTicket)
	if lastTicket != nil {
		newTicket.Position = lastTicket.Position + 1000
	}
	defer cancel()
	err2 := a.repo.CreateTicket(ctx, newTicket)
	if err2 != nil {
		log.Print(err2)
	}
}

func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	ticket := &Ticket{}

	if err := json.NewDecoder(r.Body).Decode(ticket); err != nil {
		log.Print(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := a.repo.UpdateTicket(ctx, *ticket)
	if err != nil {
		log.Print(err)
	}
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := a.repo.GetTickets(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
func (a *API) GetCalendar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := a.repo.GetTickets(ctx)
	if err != nil {
		log.Print(err)
	}
	if err := json.NewEncoder(w).Encode(result /*.ToCalendar()*/); err != nil {
		return
	}
}
