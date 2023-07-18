package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (a *API) GetMonthlyTicketsCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := a.repo.GetMonthlyTicketsCount(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
