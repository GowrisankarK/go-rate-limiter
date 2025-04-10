package main

import "fmt"
import "time"
import "github.com/GowrisankarK/go-rate-limiter/algorithm"

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

func validateSlidingWindowRateLimiter() {
	silidingWindow:=algorithm.InitialiseSlidingWindow();
	fmt.Println(fmt.Sprintf("The silidingWindow is initialised for %d request count per %d seconds", 
	silidingWindow.MaxCount, silidingWindow.Duration/1000));
	for i:=1;i<=5;i++ {
		if silidingWindow.IsRequestAllowed() {
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
	for i:=1;i<=5;i++ {
		if leakyBucket.IsRequestAllowed(int32(i)) {
			fmt.Println(fmt.Sprintf("The Request %d is allowed", i));
		} else {
			fmt.Println(fmt.Sprintf("The Request %d is not allowed", i));
		}
		time.Sleep(time.Duration(i) * time.Second);
	}
}
func main() {
	// validateFixedWindowRateLimiter();
	// validateSlidingWindowRateLimiter();
	// validateTokenBucketRateLimiter();
	validateLeakyBucketRateLimiter();
}