package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hibiken/asynq"
	"github.com/kyamalabs/users/internal/cache"
	"github.com/rs/zerolog/log"
	"github.com/wealdtech/go-ens/v3"
)

const (
	TaskCacheENSName       = "task:cache_ens_name"
	ensNameCacheKeyPrefix  = "ens-name"
	ensNameCacheExpiration = 7 * 24 * time.Hour
)

type PayloadCacheEnsName struct {
	WalletAddress string `json:"wallet_address"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCacheEnsName(ctx context.Context, payload *PayloadCacheEnsName, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskCacheENSName, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func getCacheKey(walletAddress string) string {
	return fmt.Sprintf("%s:%s", ensNameCacheKeyPrefix, walletAddress)
}

func (processor *RedisTaskProcessor) ProcessTaskCacheEnsName(ctx context.Context, task *asynq.Task) error {
	var payload PayloadCacheEnsName
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
	}

	client, err := ethclient.Dial(processor.config.EthereumRPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum blockchain at %s: %w", processor.config.EthereumRPCURL, err)
	}

	ensName, err := ens.ReverseResolve(client, common.HexToAddress(payload.WalletAddress))
	if err != nil {
		log.Info().Err(err).Str("wallet_address", payload.WalletAddress).Msg("could not resolve address into an ENS name")
	}

	err = processor.cache.Set(ctx, getCacheKey(payload.WalletAddress), ensName, ensNameCacheExpiration)
	if err != nil {
		return fmt.Errorf("could not store ens name in cache: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("processed task")

	return nil
}

func GetCachedENSName(ctx context.Context, c cache.Cache, walletAddress string) (string, error) {
	res, err := c.Get(ctx, getCacheKey(walletAddress))
	if err != nil {
		return "", fmt.Errorf("could not fetch ens name from cache: %w", err)
	}
	if res == nil {
		return "", cache.Nil
	}

	cachedENSName, ok := res.(string)
	if !ok {
		return "", errors.New("could not cast ens name to string")
	}

	return cachedENSName, nil
}
