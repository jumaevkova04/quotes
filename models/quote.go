package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Quote ...
type Quote struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Quote     string    `json:"quote"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// Quotes - Quotes + mutex
type Quotes struct {
	Quotes map[string]Quote
	sync.Mutex
}

// NewQuotes - create an instance for Quotes
func NewQuotes() *Quotes {
	return &Quotes{
		Quotes: make(map[string]Quote),
	}
}

// FillRandomlyQuotes - Fills randomly Quotes
func (q *Quotes) FillRandomlyQuotes() {
	var (
		craatedDate = time.Now()
		quoteID     = ""
	)

	q.Lock()
	defer q.Unlock()

	for i := 1; i <= 5; i++ {
		quoteID = uuid.New().String()
		craatedDate = craatedDate.Add(-20 * time.Minute)
		q.Quotes[quoteID] = Quote{
			ID:        quoteID,
			Author:    fmt.Sprintf("Author %v", i),
			Quote:     fmt.Sprintf("Quote %v", i),
			Category:  fmt.Sprintf("Category %v", i),
			CreatedAt: craatedDate,
		}
	}
}

// Add - add new quote
func (q *Quotes) Add(quote *Quote) {
	q.Lock()
	defer q.Unlock()

	quote.ID = uuid.New().String()
	q.Quotes[quote.ID] = *quote
}

// GetAll - get all quotes
func (q *Quotes) GetAll() []Quote {
	q.Lock()
	defer q.Unlock()
	var quotes []Quote
	for _, quote := range q.Quotes {
		quotes = append(quotes, quote)
	}
	return quotes
}

// GetByCategory - get quotes by category
func (q *Quotes) GetByCategory(category string) []Quote {
	q.Lock()
	defer q.Unlock()
	var quotes []Quote
	for _, quote := range q.Quotes {
		if quote.Category == category {
			quotes = append(quotes, quote)
		}
	}
	return quotes
}

// GetRandom - get random quote
func (q *Quotes) GetRandom() *Quote {
	q.Lock()
	defer q.Unlock()

	var (
		randNumber = rand.Intn(len(q.Quotes))
		counter    = 1
		quote      = Quote{}
	)

	if randNumber == 0 {
		return nil
	}

	for _, quote = range q.Quotes {
		counter++
		if counter == randNumber {
			return &quote
		}
	}
	return &quote
}

// Edit - edit quote by id
func (q *Quotes) Edit(quote *Quote) bool {
	q.Lock()
	defer q.Unlock()

	_, exists := q.Quotes[quote.ID]
	if exists {
		q.Quotes[quote.ID] = *quote
		return true
	}
	return false
}

// Delete - delete quote from quotes by id
func (q *Quotes) Delete(quoteID string) bool {
	q.Lock()
	defer q.Unlock()

	_, exists := q.Quotes[quoteID]
	if exists {
		delete(q.Quotes, quoteID)
		return true
	}
	return false
}
