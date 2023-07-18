package ticket

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	StatusNew Status = iota
	StatusInProgress
	StatusTesting
	StatusCompleted
)

type Priority int

const (
	PriorityVeryLow Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	Position    int                `bson:"position" json:"position"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description"`
	Status      Status             `bson:"status" json:"status" validate:"required"`
	Priority    Priority           `bson:"priority" json:"priority" validate:"required"`
	Tasks       []Task             `bson:"tasks" json:"tasks"`
	CompletedAt time.Time          `bson:"completed_at" json:"completedAt"`
}

type TicketByStatusCount struct {
	ID    Status `bson:"_id" json:"status"`
	Count int    `bson:"count" json:"count"`
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	Task      string             `bson:"task" json:"task" validate:"required"`
	Tag       string             `bson:"tag" json:"tag" validate:"required"`
	StartTime time.Time          `bson:"start_time" json:"startTime" validate:"required"`
	EndTime   time.Time          `bson:"end_time" json:"endTime" validate:"required"`
}

type Calendar struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Tickets []*Ticket
type Tasks []*Task

func (td Tasks) ToCalendar() (cd []Calendar) {
	for _, value := range td {
		cd = append(cd, Calendar{
			Title: value.Task,
			Start: value.StartTime,
			End:   value.EndTime,
		})
	}
	return
}

func (t *Ticket) AddCompletedAtIfCompleted() {
	if t.Status == StatusCompleted {
		t.CompletedAt = time.Now()
	}
}
