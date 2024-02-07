package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/kyamalabs/users/internal/cache"
	"github.com/kyamalabs/users/internal/util"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	QueueDefault  = "default"
	QueueCritical = "critical"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskCacheEnsName(context.Context, *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	config util.Config
	cache  cache.Cache
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskCacheENSName, processor.ProcessTaskCacheEnsName)

	return processor.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisConnOpt, config util.Config, cache cache.Cache) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueDefault:  3,
			QueueCritical: 7,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process task failed")
		}),
		Logger: logger,
	})

	return &RedisTaskProcessor{
		server: server,
		config: config,
		cache:  cache,
	}
}
