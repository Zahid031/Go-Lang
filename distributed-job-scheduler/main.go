package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
    "errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func createJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			JobType string          `json:"job_type"`
			Payload json.RawMessage `json:"payload"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var id int64
		err = pool.QueryRow(
			r.Context(),
			"INSERT INTO jobs (job_type, payload) VALUES ($1, $2) RETURNING id",
			input.JobType, input.Payload,
		).Scan(&id)
		if err != nil {
			http.Error(w, "failed to create job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     id,
			"status": "pending",
		})
	}
}
func getJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
			return
		}

		var jobType, status string
		var payload json.RawMessage
		var createdAt time.Time

		err = pool.QueryRow(
			r.Context(),
			"SELECT job_type, payload, status, created_at FROM jobs WHERE id = $1",
			id,
		).Scan(&jobType, &payload, &status, &createdAt)

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "failed to fetch job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         id,
			"job_type":   jobType,
			"payload":    payload,
			"status":     status,
			"created_at": createdAt,
		})
	}
}
func main() {
	ctx := context.Background()

	connString := "postgres://postgres:postgres@localhost:5432/job-scheduler"

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("POST /jobs", createJobHandler(pool))
	mux.HandleFunc("GET /jobs/{id}", getJobHandler(pool))

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}