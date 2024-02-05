package worker

import (
	"os"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/kyamalabs/users/internal/cache"
	"github.com/kyamalabs/users/internal/util"
	"github.com/rs/zerolog/log"
)

var (
	testTaskDistributor TaskDistributor
	testRedisCache      cache.Cache
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}

	redisOpt, err := asynq.ParseRedisURI(config.RedisConnURL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse redis URL")
	}

	testRedisCache, err = cache.NewRedisCache(config.RedisConnURL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create new redis cache")
	}

	testTaskDistributor = NewRedisTaskDistributor(redisOpt)

	go runTestTaskProcessor(config, redisOpt, testRedisCache)

	os.Exit(m.Run())
}

func runTestTaskProcessor(config util.Config, redisOpt asynq.RedisConnOpt, redisCache cache.Cache) {
	taskProcessor := NewRedisTaskProcessor(redisOpt, config, redisCache)

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("could not start test task processor")
	}

	log.Info().Msg("started test task processor")
}
