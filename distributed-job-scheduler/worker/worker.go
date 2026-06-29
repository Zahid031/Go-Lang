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

const NumWorkers = 3

type Worker struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func New(pool *pgxpool.Pool, redisClient *redis.Client) *Worker {
	return &Worker{pool: pool, redis: redisClient}
}

func (w *Worker) Run(ctx context.Context) {
	log.Printf("starting %d workers, waiting for jobs...", NumWorkers)

	for i := 1; i <= NumWorkers; i++ {
		go w.loop(ctx, i)
	}

	// Block forever so Run doesn't return immediately while goroutines
	// keep working in the background. We'll replace this with a proper
	// shutdown signal later.
	select {}
}

func (w *Worker) loop(ctx context.Context, workerID int) {
	for {
		result, err := w.redis.BLPop(ctx, 5*time.Second, JobQueueKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("[worker %d] error popping from queue: %v", workerID, err)
			continue
		}

		jobID := result[1]
		w.processJob(ctx, workerID, jobID)
	}
}

func (w *Worker) processJob(ctx context.Context, workerID int, jobIDStr string) {
	id, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		log.Printf("[worker %d] invalid job id from queue: %s", workerID, jobIDStr)
		return
	}

	var jobType string
	err = w.pool.QueryRow(ctx,
		"SELECT job_type FROM jobs WHERE id = $1",
		id,
	).Scan(&jobType)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("[worker %d] job %d not found in db, skipping", workerID, id)
		return
	}
	if err != nil {
		log.Printf("[worker %d] failed to fetch job %d: %v", workerID, id, err)
		return
	}

	log.Printf("[worker %d] starting job %d (type=%s)", workerID, id, jobType)

	_, err = w.pool.Exec(ctx,
		"UPDATE jobs SET status = 'running' WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("[worker %d] failed to mark job %d as running: %v", workerID, id, err)
		return
	}

	// Simulated work — real execution logic (send-email, resize-image, etc.) comes later.
	time.Sleep(2 * time.Second)

	_, err = w.pool.Exec(ctx,
		"UPDATE jobs SET status = 'completed' WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("[worker %d] failed to mark job %d as completed: %v", workerID, id, err)
		return
	}

	log.Printf("[worker %d] finished job %d", workerID, id)
}