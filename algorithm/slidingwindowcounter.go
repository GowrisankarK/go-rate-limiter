package algorithm

import "fmt"
import "time"

type SlidingWindowCounter struct {
	MaxCount int32 
	Duration int64 // milliseconds
	Buckets map[int64]int
}

func(slidingWindowCounter *SlidingWindowCounter) cleanupOldRequests(timeStamp int64) {
	fmt.Println("Clean up started");
	for ts,_:=range(slidingWindowCounter.Buckets) {
		if ts<=timeStamp {
			delete(slidingWindowCounter.Buckets, ts);
		}
	}
}

func(slidingWindowCounter *SlidingWindowCounter) IsRequestAllowed() bool {
	now := time.Now().UnixMilli();
	timeFrame:=now-slidingWindowCounter.Duration;
	defer slidingWindowCounter.cleanupOldRequests(timeFrame);
	var reqCount int32 = 0;
	for ts,_ := range(slidingWindowCounter.Buckets) {
		if ts > timeFrame {
			reqCount++;
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
	return slidingWindowCounter;
}