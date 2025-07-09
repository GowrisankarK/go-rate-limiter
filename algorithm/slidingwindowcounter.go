package algorithm

import "fmt"
import "time"
import "sync"

type SlidingWindowCounter struct {
	MaxCount	int32 
	Duration 	int64 // milliseconds
	Buckets  	map[int64]int
	mutex       sync.RWMutex
}

func (slidingWindowCounter *SlidingWindowCounter) cleanupOldRequestsLoop() {
	ticker := time.NewTicker(1 * time.Second); // run cleanup every second
	for range ticker.C {
		cutoff := (time.Now().UnixMilli() - slidingWindowCounter.Duration*2); // cutoff of the provious window

		// Step 1: Quickly collect up to 2 expired timestamps
		var oldTimestamps []int64;
		slidingWindowCounter.mutex.RLock();
		count := 0;
		for ts := range slidingWindowCounter.Buckets {
			if ts <= cutoff {
				oldTimestamps = append(oldTimestamps, ts);
				count++;
				if count >= 2 { // clean only 2 entries at a time
					break;
				}
			}
		}
		slidingWindowCounter.mutex.RUnlock();

		// Step 2: Delete them outside if needed
		if len(oldTimestamps) > 0 {
			slidingWindowCounter.mutex.Lock();
			for _, ts := range oldTimestamps {
				delete(slidingWindowCounter.Buckets, ts)
			}
			slidingWindowCounter.mutex.Unlock();
		}

		if len(oldTimestamps) > 0 {
			fmt.Println("Cleaned up", len(oldTimestamps), "old entries");
		}
	}
}

func(slidingWindowCounter *SlidingWindowCounter) IsRequestAllowed() bool {
	now := time.Now().UnixMilli();
	timeFrame:=now-slidingWindowCounter.Duration;
	
	var reqCount int32 = 0;
	slidingWindowCounter.mutex.Lock();
	defer slidingWindowCounter.mutex.Unlock();
	for ts,count := range(slidingWindowCounter.Buckets) {
		if ts > timeFrame {
			reqCount+=int32(count);
		}
	}
	if reqCount < slidingWindowCounter.MaxCount {
		slidingWindowCounter.Buckets[now]++;
		return true;
	}
	return false;
}

func InitialiseSlidingWindowCounter() SlidingWindowCounter {
	slidingWindowCounter := SlidingWindowCounter{MaxCount: 2, Duration: 10000, Buckets: map[int64]int{}};
	go slidingWindowCounter.cleanupOldRequestsLoop() // background cleanup starts
	return slidingWindowCounter;
}