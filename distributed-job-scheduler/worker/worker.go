package worker

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const JobQueueKey = "jobs:queue"

type Worker struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func New(pool *pgxpool.Pool, redisClient *redis.Client) *Worker {
	return &Worker{pool: pool, redis: redisClient}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("worker started, waiting for jobs...")

	for {
		result, err := w.redis.BLPop(ctx, 5*time.Second, JobQueueKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("error popping from queue: %v", err)
			continue
		}

		jobID := result[1]
		w.processJob(ctx, jobID)
	}
}

func (w *Worker) processJob(ctx context.Context, jobIDStr string) {
	id, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		log.Printf("invalid job id from queue: %s", jobIDStr)
		return
	}

	var jobType string
	err = w.pool.QueryRow(ctx,
		"SELECT job_type FROM jobs WHERE id = $1",
		id,
	).Scan(&jobType)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("job %d not found in db, skipping", id)
		return
	}
	if err != nil {
		log.Printf("failed to fetch job %d: %v", id, err)
		return
	}

	log.Printf("starting job %d (type=%s)", id, jobType)

	_, err = w.pool.Exec(ctx,
		"UPDATE jobs SET status = 'running' WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("failed to mark job %d as running: %v", id, err)
		return
	}

	// Simulated work — real execution logic (send-email, resize-image, etc.) comes later.
	time.Sleep(2 * time.Second)

	_, err = w.pool.Exec(ctx,
		"UPDATE jobs SET status = 'completed' WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("failed to mark job %d as completed: %v", id, err)
		return
	}

	log.Printf("finished job %d", id)
}