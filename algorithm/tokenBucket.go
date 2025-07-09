package algorithm

import "fmt"
import "time"
import "sync"

type TokenBucket struct {
	AvailableTokenCount int32
	ReFillTokenCount int32
	MaxCount int32
	RefillRate int64
	ticker     *time.Ticker
	stopChannel   chan struct{}
	mutex         sync.Mutex
}

func(tokenBucket *TokenBucket) IsRequestAllowed() bool {
	tokenBucket.mutex.Lock();
	defer tokenBucket.mutex.Unlock()
	fmt.Println(fmt.Sprintf("The Available Tokens %d", tokenBucket.AvailableTokenCount));
	if tokenBucket.AvailableTokenCount > 0 {
		tokenBucket.AvailableTokenCount--;
		return true;
	}
	return false;
}

func (tokenBucket *TokenBucket) refillTokens() {
	fmt.Println("Refill token is called")
	tokenBucket.ticker = time.NewTicker(time.Duration(tokenBucket.RefillRate) * time.Millisecond)
	tokenBucket.stopChannel = make(chan struct{})

	go func() {
		for {
			select {
			case <-tokenBucket.ticker.C:
				tokenBucket.mutex.Lock();
				if tokenBucket.AvailableTokenCount < tokenBucket.MaxCount {
					tokenBucket.AvailableTokenCount += tokenBucket.ReFillTokenCount
					if tokenBucket.AvailableTokenCount > tokenBucket.MaxCount {
						tokenBucket.AvailableTokenCount = tokenBucket.MaxCount
					}
					fmt.Printf("Refilled %d token(s), current token count: %d\n", tokenBucket.ReFillTokenCount, tokenBucket.AvailableTokenCount)
				}
				tokenBucket.mutex.Unlock()
			case <-tokenBucket.stopChannel:
				fmt.Println("Token Refill Stopped")
				tokenBucket.ticker.Stop()
				return
			}
		}
	}()
}

func (tokenBucket *TokenBucket) StopRefill() {
	close(tokenBucket.stopChannel)
}

func InitialiseTokenBucket() *TokenBucket {
	tokenBucket := &TokenBucket{AvailableTokenCount: 3, ReFillTokenCount: 3, MaxCount: 10, RefillRate: 1000};
	tokenBucket.refillTokens();
	return tokenBucket;
}