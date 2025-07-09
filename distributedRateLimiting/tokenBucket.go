package distributedRateLimiting

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBucket struct {
	RedisClient      *redis.Client
	ReFillTokenCount int
	MaxCount         int
	RefillRate       time.Duration

	mu           sync.Mutex
	tickers      map[string]*time.Ticker
	stopChannels map[string]chan struct{}
}

var ctx = context.Background()

func NewRedisTokenBucket(redisClient *redis.Client, refillRate time.Duration, refillCount int, maxCount int, clientId string) *RedisTokenBucket {
	redisTokenBucket := &RedisTokenBucket{
		RedisClient:      redisClient,
		ReFillTokenCount: refillCount,
		MaxCount:         maxCount,
		RefillRate:       refillRate,
		tickers:          make(map[string]*time.Ticker),
		stopChannels:     make(map[string]chan struct{}),
	}
	redisTokenBucket.StartRefill(clientId);
	return redisTokenBucket;
}

// Start token refill for a specific clientID
func (tokenBucket *RedisTokenBucket) StartRefill(clientID string) {
	tokenBucket.mu.Lock()
	if _, exists := tokenBucket.stopChannels[clientID]; exists {
		tokenBucket.mu.Unlock()
		return // already started
	}

	ticker := time.NewTicker(tokenBucket.RefillRate)
	stopChan := make(chan struct{})
	tokenBucket.tickers[clientID] = ticker
	tokenBucket.stopChannels[clientID] = stopChan
	tokenBucket.mu.Unlock()

	// Initialize Redis keys if not already present
	keyTokens := fmt.Sprintf("tokenbucket:%s:tokens", clientID)
	keyTimestamp := fmt.Sprintf("tokenbucket:%s:timestamp", clientID)

	// Use SetNX (Set if Not Exists)
	now := time.Now().UnixMilli()
	tokenBucket.RedisClient.SetNX(ctx, keyTokens, tokenBucket.MaxCount, 0)
	tokenBucket.RedisClient.SetNX(ctx, keyTimestamp, now, 0)

	go func() {
		for {
			select {
			case <-ticker.C:
				tokenBucket.refill(clientID)
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop token refill for a specific clientID
func (tokenBucket *RedisTokenBucket) StopRefill(clientID string) {
	tokenBucket.mu.Lock()
	defer tokenBucket.mu.Unlock()

	if stopChan, exists := tokenBucket.stopChannels[clientID]; exists {
		close(stopChan)
		delete(tokenBucket.stopChannels, clientID)
	}

	if ticker, exists := tokenBucket.tickers[clientID]; exists {
		ticker.Stop()
		delete(tokenBucket.tickers, clientID)
	}
}

// Main refill logic (shared by all clientIDs)
func (tokenBucket *RedisTokenBucket) refill(clientID string) {
	keyTokens := fmt.Sprintf("tokenbucket:%s:tokens", clientID)
	keyTimestamp := fmt.Sprintf("tokenbucket:%s:timestamp", clientID)
	now := time.Now().UnixMilli()

	pipe := tokenBucket.RedisClient.TxPipeline()
	tokensCmd := pipe.Get(ctx, keyTokens)
	timestampCmd := pipe.Get(ctx, keyTimestamp)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		fmt.Println("Refill Exec error:", err)
		return
	}

	tokenCount := 0
	lastRefill := now

	if t, err := tokensCmd.Int(); err == nil {
		tokenCount = t
	}

	if ts, err := timestampCmd.Int64(); err == nil {
		lastRefill = ts
	}

	elapsed := now - lastRefill
	refillTokens := int(elapsed / tokenBucket.RefillRate.Milliseconds()) * tokenBucket.ReFillTokenCount

	if refillTokens > 0 {
		tokenCount = min(tokenCount+refillTokens, tokenBucket.MaxCount)

		pipe := tokenBucket.RedisClient.TxPipeline()
		pipe.Set(ctx, keyTokens, tokenCount, 0)
		pipe.Set(ctx, keyTimestamp, now, 0)
		_, err := pipe.Exec(ctx)
		if err != nil {
			fmt.Println("Refill Set error:", err)
		}
	}
}

// Check if a request is allowed for the given clientID
func (tokenBucket *RedisTokenBucket) IsRequestAllowed(clientID string) (bool, error) {
	keyTokens := fmt.Sprintf("tokenbucket:%s:tokens", clientID)

	pipe := tokenBucket.RedisClient.TxPipeline()
	tokensCmd := pipe.Get(ctx, keyTokens)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}

	tokenCount := 0
	if t, err := tokensCmd.Int(); err == nil {
		tokenCount = t
	}

	if tokenCount <= 0 {
		return false, nil
	}

	tokenCount--

	pipe = tokenBucket.RedisClient.TxPipeline()
	pipe.Set(ctx, keyTokens, tokenCount, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
