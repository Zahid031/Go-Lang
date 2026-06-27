package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zahid031/distributed-job-scheduler/models"
)

// func CreateJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
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


func CreateJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

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
			"INSERT INTO jobs (job_type, payload, user_id) VALUES ($1, $2, $3) RETURNING id",
			input.JobType, input.Payload, userID,
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


func GetJobHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
			return
		}

		var j models.Job
		err = pool.QueryRow(
			r.Context(),
			"SELECT id, job_type, payload, status, created_at FROM jobs WHERE id = $1 AND user_id = $2",
			id, userID,
		).Scan(&j.ID, &j.JobType, &j.Payload, &j.Status, &j.CreatedAt)

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
		json.NewEncoder(w).Encode(j)
	}
}

func ListJobsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := pool.Query(
			r.Context(),
			"SELECT id, job_type, payload, status, created_at FROM jobs WHERE user_id = $1",
			userID,
		)
		if err != nil {
			http.Error(w, "failed to fetch jobs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var jobs []models.Job
		for rows.Next() {
			var j models.Job
			err := rows.Scan(&j.ID, &j.JobType, &j.Payload, &j.Status, &j.CreatedAt)
			if err != nil {
				http.Error(w, "failed to scan job", http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, j)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "error iterating over jobs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jobs)
	}
}