package dashboard

type yearMonth struct {
	Year  int `bson:"year" json:"year"`
	Month int `bson:"month" json:"month"`
}
type MonthlyTicket struct {
	ID              yearMonth `bson:"_id" json:"yearMonth"`
	NumberOfTickets int       `bson:"numberOfTickets" json:"numberOfTickets"`
}

type MonthlyTickets []*MonthlyTicket
