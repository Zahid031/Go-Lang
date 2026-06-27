package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/zahid031/distributed-job-scheduler/handlers"

)
type Job struct {
	ID        int64           `json:"id"`
	JobType   string          `json:"job_type"`
	Payload   json.RawMessage `json:"payload"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// func createJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var input struct {
// 			JobType string          `json:"job_type"`
// 			Payload json.RawMessage `json:"payload"`
// 		}

// 		err := json.NewDecoder(r.Body).Decode(&input)
// 		if err != nil {
// 			http.Error(w, "invalid request body", http.StatusBadRequest)
// 			return
// 		}

// 		var id int64
// 		err = pool.QueryRow(
// 			r.Context(),
// 			"INSERT INTO jobs (job_type, payload) VALUES ($1, $2) RETURNING id",
// 			input.JobType, input.Payload,
// 		).Scan(&id)
// 		if err != nil {
// 			http.Error(w, "failed to create job", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(map[string]any{
// 			"id":     id,
// 			"status": "pending",
// 		})
// 	}
// }


// func getJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		idStr := chi.URLParam(r, "id")
// 		id, err := strconv.ParseInt(idStr, 10, 64)
// 		if err != nil {
// 			http.Error(w, "invalid job id", http.StatusBadRequest)
// 			return
// 		}

// 		var j Job
// 		err = pool.QueryRow(
// 			r.Context(),
// 			"SELECT id, job_type, payload, status, created_at FROM jobs WHERE id = $1",
// 			id,
// 		).Scan(&j.ID, &j.JobType, &j.Payload, &j.Status, &j.CreatedAt)

// 		if errors.Is(err, pgx.ErrNoRows) {
// 			http.Error(w, "job not found", http.StatusNotFound)
// 			return
// 		}
// 		if err != nil {
// 			http.Error(w, "failed to fetch job", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(j)
// 	}
// }

// func listJobsHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rows, err := pool.Query(
// 			r.Context(),
// 			"SELECT id, job_type, payload, status, created_at FROM jobs",
// 		)
// 		if err != nil {
// 			http.Error(w, "failed to fetch jobs", http.StatusInternalServerError)
// 			return 
// 		}
// 		defer rows.Close()
// 		var jobs []Job 
// 		for rows.Next() {
// 			var j Job
// 			err := rows.Scan(&j.ID, &j.JobType, &j.Payload, &j.Status, &j.CreatedAt)
// 			if err != nil {
// 				http.Error(w, "failed to scan job", http.StatusInternalServerError)
// 				return
// 			}
// 			jobs = append(jobs, j)
// 		}
// 		if err := rows.Err(); err != nil {
// 			http.Error(w, "error iterating over jobs", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(jobs)
// 	}
// } 

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
func main() {
	ctx := context.Background()
	if err:= godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// connString := "postgres://postgres:postgres@localhost:5432/job-scheduler"
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")	
	}

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
	r := chi.NewRouter()
	r.Use(loggingMiddleware)
	r.Get("/health", healthHandler)
	r.Post("/jobs", handlers.CreateJobHandler(pool))
	r.Get("/jobs/{id}", handlers.GetJobHandler(pool))
	r.Get("/jobs", handlers.ListJobsHandler(pool))

	// mux := http.NewServeMux()
	// mux.HandleFunc("/health", healthHandler)
	// log.Println("test log...")
	// mux.HandleFunc("POST /jobs", createJobHandler(pool))
	// mux.HandleFunc("GET /jobs/{id}", getJobHandler(pool))
	// mux.HandleFunc("GET /jobs", listJobsHandler(pool))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}