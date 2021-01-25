package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/jumaevkova04/quotes/handlers"
)

func main() {
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("------------Start Logging------------")

	router := httprouter.New()

	router.POST("/quotes", handlers.AddQuotes)                   // 1. Create a quote with following fields: author, quote, category.
	router.GET("/quotes", handlers.GetAllQuotes)                 // 4. Get all quotes.
	router.GET("/randomquotes", handlers.GetRandomQuote)         // 6. Get a random quote.
	router.GET("/quotes/:category", handlers.GetQuoteByCategory) // 5. Get all quotes by category.
	router.PUT("/quotes/:id", handlers.EditQuoteByID)            // 2. Edit/Change quote: author, quote, category.
	router.DELETE("/quotes/:id", handlers.DeleteQuoteByID)       // 3. Delete a quote by ID.

	// a worker that wakes up every 5 minutes and deletes quotes
	go worker(time.Minute*5, handlers.DeleteOldQuotes)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":3939", router))

	time.Sleep(time.Minute * 2)
}

func worker(d time.Duration, f func()) {
	var reentrancyFlag int64
	for range time.Tick(d) {
		if atomic.CompareAndSwapInt64(&reentrancyFlag, 0, 1) {
			defer atomic.StoreInt64(&reentrancyFlag, 0)
		} else {
			return
		}
		f()
	}
}
