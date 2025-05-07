package main

import "fmt"
import "time"
import "github.com/GowrisankarK/go-rate-limiter/algorithm"
import "github.com/GowrisankarK/go-rate-limiter/distributedRateLimiting"
import "github.com/redis/go-redis/v9"

func validateFixedWindowRateLimiter() {
	fixedWindow:=algorithm.InitialiseFixedWindow();
	fmt.Println(fmt.Sprintf("The fixedWindow is initialised for %d request count per %d seconds and start time %d", 
	fixedWindow.MaxCount, fixedWindow.Duration/1000, fixedWindow.StartTimestamp));
	for i:=1;i<=100;i++ {
		if fixedWindow.IsRequestAllowed() {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(5 * time.Second);
	}
}

func validateSlidingWindowLogRateLimiter() {
	silidingWindowLog:=algorithm.InitialiseSlidingWindow();
	fmt.Println(fmt.Sprintf("The silidingWindow is initialised for %d request count per %d seconds", 
	silidingWindowLog.MaxCount, silidingWindowLog.Duration/1000));
	for i:=1;i<=20;i++ {
		if silidingWindowLog.IsRequestAllowed() {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(time.Duration(i) * time.Second);
	}
}

func validateSlidingWindowCounterRateLimiter() {
	silidingWindowCounter:=algorithm.InitialiseSlidingWindowCounter();
	fmt.Println(fmt.Sprintf("The silidingWindow is initialised for %d request count per %d seconds", 
	silidingWindowCounter.MaxCount, silidingWindowCounter.Duration/1000));
	for i:=1;i<=20;i++ {
		if silidingWindowCounter.IsRequestAllowed() {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(time.Duration(i) * time.Second);
	}
}

func validateTokenBucketRateLimiter() {
	tokenBucket:=algorithm.InitialiseTokenBucket();
	fmt.Println(fmt.Sprintf("The tokenBucket is initialised for %d request count per %d seconds", 
	tokenBucket.ReFillTokenCount, tokenBucket.RefillRate/1000));
	for i:=1;i<=5;i++ {
		if tokenBucket.IsRequestAllowed() {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(time.Duration(i) * time.Second);
	}
}

func validateLeakyBucketRateLimiter() {
	leakyBucket:=algorithm.InitialiseLeakyBucket();
	fmt.Println(fmt.Sprintf("The leakyBucket is initialised to process request per %d seconds with max request queue size of %d", 
	leakyBucket.OutFillRate/1000, leakyBucket.MaxTokenCount));
	for i:=1;i<=10;i++ {
		if leakyBucket.IsRequestAllowed(int32(i)) {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(time.Duration(i) * time.Second);
	}
}

func validateRedisTokenBucket() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	tokenBucket := distributedRateLimiting.NewRedisTokenBucket(client, 5000*time.Millisecond, 10, 20)

	clientID := "b21295cb-5e59-4730-83b3-c39b57afd4a2"
	tokenBucket.StartRefill(clientID);

	fmt.Printf("The RedisTokenBucket is initialised for %d refill tokens per %d seconds, max tokens = %d\n",
		tokenBucket.ReFillTokenCount, tokenBucket.RefillRate/time.Second, tokenBucket.MaxCount)

	for i := 1; i <= 20; i++ {
		allowed, err := tokenBucket.IsRequestAllowed(clientID)
		if err != nil {
			fmt.Printf("Request %d failed with error: %v\n", i, err)
			continue
		}

		if allowed {
			fmt.Printf("Request %d is allowed\n", i)
		} else {
			fmt.Printf("Request %d is not allowed\n", i)
		}

		time.Sleep(1 * time.Second)
	}
	tokenBucket.StopRefill(clientID);
}


func main() {
	// validateFixedWindowRateLimiter();
	// validateSlidingWindowLogRateLimiter();
	// validateSlidingWindowCounterRateLimiter();
	// validateTokenBucketRateLimiter();
	// validateLeakyBucketRateLimiter();
	validateRedisTokenBucket();
}