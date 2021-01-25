package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/jumaevkova04/quotes/helpers"
	"github.com/jumaevkova04/quotes/models"
)

var (
	quotes = models.NewQuotes()
)

func init() {
	quotes.FillRandomlyQuotes()
}

// AddQuotes - add new quote
func AddQuotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		quote    models.Quote
		response = models.Response{
			Code: 200,
		}
		err error
	)
	defer response.Send(w, r)

	err = json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		response.Code = http.StatusBadRequest
		log.Println("handlers:AddQuotes parcing request failed:", err)
		return
	}
	quotes.Add(&quote)
	response.Payload = quote
}

// GetAllQuotes - get all quotes
func GetAllQuotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		response = models.Response{
			Code: 200,
		}
	)
	defer response.Send(w, r)

	response.Payload = quotes.GetAll()
}

// GetQuoteByCategory - get quotes by category
func GetQuoteByCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		response = models.Response{
			Code: 200,
		}
	)
	defer response.Send(w, r)

	category := ps.ByName("category")
	category = helpers.RemoveAllSpaces(category)
	if len(category) < 2 {
		response.Code = http.StatusBadRequest
		return
	}

	response.Payload = quotes.GetByCategory(category)
}

// GetRandomQuote - get some random quote
func GetRandomQuote(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		response = models.Response{
			Code: 200,
		}
	)
	defer response.Send(w, r)

	response.Payload = quotes.GetRandom()
}

// EditQuoteByID - edit quote by id
func EditQuoteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		quote    models.Quote
		response = models.Response{
			Code: 200,
		}
		err error
	)
	defer response.Send(w, r)

	quoteID := ps.ByName("id")
	if !helpers.IsValidUUID(quoteID) {
		response.Code = http.StatusBadRequest
		response.Message = "wrong quote id"
		log.Println("handlers:EditQuote got wrong quote id(uuid) error:", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		response.Code = http.StatusBadRequest
		log.Println("handlers:EditQuote parcing request failed:", err)
		return
	}

	if !quotes.Edit(&quote) {
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("quote with id[%s] not exists", quoteID)
		log.Println("handlers:EditQuote wrong quote id(uuid) error: not exists")
		return
	}

	response.Payload = quotes.GetAll()
}

// DeleteQuoteByID - delete quote by id
func DeleteQuoteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		response = models.Response{
			Code: 200,
		}
		err error
	)
	defer response.Send(w, r)

	quoteID := ps.ByName("id")
	if !helpers.IsValidUUID(quoteID) {
		response.Code = http.StatusBadRequest
		response.Message = "wrong quote id"
		log.Println("handlers:DeleteQuoteByID got wrong quote id(uuid) error:", err)
		return
	}

	if !quotes.Delete(quoteID) {
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("quote with id[%s] not exists", quoteID)
		log.Println("handlers:DeleteQuoteByID wrong quote id(uuid) error: not exists")
		return
	}

	response.Payload = quotes.GetAll()
}

// DeleteOldQuotes - deletes quotes that were created more than 1 hour ago.
func DeleteOldQuotes() {
	for _, quote := range quotes.GetAll() {
		if helpers.IsTimePassed(time.Now().Add(-time.Hour), quote.CreatedAt) {
			quotes.Delete(quote.ID)
		}
	}
}
