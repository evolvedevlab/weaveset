package data

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Handler interface {
	Handle(context.Context, *Job) error
}

type Queuer interface {
	Enqueue(context.Context, *Job) error
	Consume(context.Context, Handler) error
}

type Job struct {
	ID        string
	URL       string
	CreatedAt time.Time
}

type RedisQueue struct {
	hostname string
	stream   string
	group    string
	client   *redis.Client
}

func NewRedisQueue(hostname, stream, group string, client *redis.Client) Queuer {
	return &RedisQueue{
		hostname: hostname,
		stream:   stream,
		group:    group,
		client:   client,
	}
}

func (q *RedisQueue) Consume(ctx context.Context, handler Handler) error {
	for {
		streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    q.group,
			Consumer: q.hostname,
			Streams:  []string{q.stream, ">"},
			Count:    10,
			Block:    0, // block forever
		}).Result()
		if err != nil {
			slog.Error("Consume error", "hostname", q.hostname, "err", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				raw, ok := msg.Values["data"].(string)
				if !ok {
					slog.Error("Consume error", "hostname", q.hostname, "err", "invalid payload data")
					continue
				}

				var job Job
				if err := json.Unmarshal([]byte(raw), &job); err != nil {
					slog.Error("Consume error", "hostname", q.hostname, "err", err)
					continue
				}

				if err := handler.Handle(ctx, &job); err != nil {
					slog.Error("Consume error", "hostname", q.hostname, "err", err)
					continue
				}

				if err := q.client.XAck(ctx, q.stream, q.group, msg.ID).Err(); err != nil {
					slog.Error("Consume error", "hostname", q.hostname, "err", err)
				}
			}
		}
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *Job) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]any{
			"data": payload,
		},
	}).Err()
	return err
}
