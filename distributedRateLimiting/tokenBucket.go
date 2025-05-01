package distributedRateLimiting

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBucket struct {
	RedisClient      *redis.Client
	ReFillTokenCount int
	MaxCount         int
	RefillRate       time.Duration // milliseconds
}

var ctx = context.Background()

func NewRedisTokenBucket(redisClient *redis.Client, refillRate time.Duration, refillCount int, maxCount int) *RedisTokenBucket {
	return &RedisTokenBucket{
		RedisClient:      redisClient,
		RefillRate:       refillRate,
		ReFillTokenCount: refillCount,
		MaxCount:         maxCount,
	}
}

func (tb *RedisTokenBucket) refillTokens(tokenCount int, lastRefill int64, now int64) (int, int64) {
	elapsed := now - lastRefill
	refillTokens := int(elapsed / tb.RefillRate.Milliseconds()) * tb.ReFillTokenCount
	if refillTokens > 0 {
		tokenCount = min(tokenCount+refillTokens, tb.MaxCount)
		lastRefill = now
	}
	return tokenCount, lastRefill
}


func (tb *RedisTokenBucket) IsRequestAllowed(clientID string) (bool, error) {
	keyTokens := fmt.Sprintf("tokenbucket:%s:tokens", clientID)
	keyTimestamp := fmt.Sprintf("tokenbucket:%s:timestamp", clientID)

	now := time.Now().UnixMilli()

	// Fetch current token count and last refill time
	pipe := tb.RedisClient.TxPipeline()
	tokensCmd := pipe.Get(ctx, keyTokens)
	timestampCmd := pipe.Get(ctx, keyTimestamp)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}

	var tokenCount int
	var lastRefill int64

	if tokens, err := tokensCmd.Int(); err == nil {
		tokenCount = tokens
	}
	if ts, err := timestampCmd.Int64(); err == nil {
		lastRefill = ts
	} else {
		lastRefill = now
	}

	// Use separated refill function
	tokenCount, lastRefill = tb.refillTokens(tokenCount, lastRefill, now)

	allowed := false
	if tokenCount > 0 {
		tokenCount--
		allowed = true
	}

	// Save back to Redis
	pipe = tb.RedisClient.TxPipeline()
	pipe.Set(ctx, keyTokens, tokenCount, 0)
	pipe.Set(ctx, keyTimestamp, lastRefill, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return allowed, nil
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
