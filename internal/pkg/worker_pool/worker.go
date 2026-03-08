package worker_pool

import (
	"context"
	"fmt"
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

type Task interface {
	Do(ctx context.Context) error
}

type Worker struct {
	task     Task
	schedule cron.Schedule
	location *time.Location
}

func NewWorker(task Task, cronTime string, location *time.Location) Worker {
	schedule, err := cron.ParseStandard(cronTime)
	if err != nil {
		panic(fmt.Sprintf("failed to parse cron: %s", err.Error()))
	}

	return Worker{
		task:     task,
		schedule: schedule,
		location: location,
	}
}

func (w Worker) Start(ctx context.Context) {
	for {
		now := time.Now().UTC()
		if w.location != nil {
			now = time.Now().In(w.location)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(w.nextScheduledTime(now).Sub(now)):
			if err := w.task.Do(ctx); err != nil {
				log.Printf("task process error: %v", err)
			}
		}
	}
}

func (w Worker) nextScheduledTime(now time.Time) time.Time {
	return w.schedule.Next(now)
}

type WorkerPool struct {
	workers []Worker
}

func NewWorkerPool(workers []Worker) WorkerPool {
	return WorkerPool{
		workers: workers,
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	for _, w := range wp.workers {
		go w.Start(ctx)
	}
}
