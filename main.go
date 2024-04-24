package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func main() {

	db, err := initStore()
	if err != nil {
		log.Fatal("Failed to initialise the store: %s", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello dev"))
	})

	mux.HandleFunc("GET /healthy", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.MarshalIndent(
			map[string]interface{}{"Status": "OK"},
			"",
			"\t",
		)
		if err != nil {
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	})

	server := http.Server{
		Addr:         ":8888",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("server started!")
	log.Fatal(server.ListenAndServe())
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func initStore() (*sql.DB, error) {

	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	var (
		db  *sql.DB
		err error
	)

	openDB := func() error {
		db, err = sql.Open("postgres", pgConnString)
		return err
	}

	err = backoff.Retry(openDB, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS message (value TEXT PRIMARY KEY)"); err != nil {
		return nil, err
	}

	return db, nil
}
