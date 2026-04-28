package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/evolvedevlab/weaveset/config"
	"github.com/evolvedevlab/weaveset/data"
)

// WARN: Not safe to use. Need synchronization...

const (
	maxInFlight = 1000
	backoffBase = time.Second * 2
)

type wMessage struct {
	*data.Job
	retries  int
	inFlight bool
	last     time.Time
}

func newWMessage(job *data.Job) *wMessage {
	return &wMessage{
		Job: job,
	}
}

type WorkerPool struct {
	n     int // number of workers
	msgch chan *wMessage
	data  map[string]*wMessage
}

func NewWorkerPool(concurrency, buffer int) Queuer {
	if concurrency == 0 {
		concurrency = 10
	}
	if buffer == 0 {
		buffer = 100
	}
	return &WorkerPool{
		n:     concurrency,
		msgch: make(chan *wMessage, buffer),
		data:  make(map[string]*wMessage),
	}
}

func (q *WorkerPool) Consume(ctx context.Context, handler data.Handler) error {
	for i := 0; i < q.n; i++ {
		go q.worker(ctx, handler)
	}
	go q.reaperLoop(ctx)

	<-ctx.Done()
	return ctx.Err()
}

func (q *WorkerPool) Enqueue(ctx context.Context, job *data.Job) error {
	if len(q.data) >= maxInFlight {
		return fmt.Errorf("queue full")
	}

	msg := newWMessage(job)
	q.data[msg.ID] = msg

	select {
	case q.msgch <- msg:
		return nil
	case <-ctx.Done():
		delete(q.data, job.ID)
		return ctx.Err()
	}
}

func (q *WorkerPool) reaperLoop(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, msg := range q.data {
				if q.shouldRetry(msg.ID) {
					q.msgch <- msg
				}
			}
		}
	}
}

func (q *WorkerPool) worker(ctx context.Context, handler data.Handler) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-q.msgch:
			msg.inFlight = true
			if err := q.handleMessage(ctx, msg, handler); err != nil {
				if errors.Is(err, errExceededRetryLimit) {
					slog.Info("job_dropped", "msg_id", msg.ID, "err", err)
				} else {
					slog.Error("job_processing_failed",
						"op", "worker_handleMessage",
						"msg_id", msg.ID,
						"err", err)
				}
			}

			msg.last = time.Now()
			msg.inFlight = false
		}
	}
}

func (q *WorkerPool) handleMessage(ctx context.Context, msg *wMessage, handler data.Handler) error {
	retries, err := q.getRetryCount(msg.Job.ID)
	if err != nil {
		return err
	}

	if retries >= config.MaxJobRetryLimit {
		return q.dropMessage(msg.ID)
	}

	err = q.processMessage(ctx, msg, handler)
	return err
}

func (q *WorkerPool) processMessage(ctx context.Context, msg *wMessage, handler data.Handler) error {
	err := handler.Handle(ctx, msg.Job)
	if err != nil {
		// increment retry count
		if err := q.incrRetryCount(msg.Job.ID); err != nil {
			return err
		}
		return err
	}

	delete(q.data, msg.ID)
	return nil
}

func (q *WorkerPool) dropMessage(msgID string) error {
	delete(q.data, msgID)
	return errExceededRetryLimit
}

func (q *WorkerPool) shouldRetry(msgID string) bool {
	msg, ok := q.data[msgID]
	if !ok {
		return false
	}

	wait := backoffBase * time.Duration(1<<msg.retries)
	return !msg.inFlight && time.Since(msg.last) >= wait
}

func (q *WorkerPool) incrRetryCount(jobID string) error {
	msg, ok := q.data[jobID]
	if !ok {
		return fmt.Errorf("message not found") // not found
	}

	msg.retries++
	return nil
}

func (q *WorkerPool) getRetryCount(jobID string) (int, error) {
	msg, ok := q.data[jobID]
	if !ok {
		return 0, nil // not found
	}

	return msg.retries, nil
}
