package algorithm

import "fmt"
import "time"
import "sync"

type SlidingWindow struct {
	MaxCount 		int32
	Duration 		int64
	RequestHistory 	[]int64 // milliseconds
	mutex         	sync.RWMutex
}

func (slidingWindow *SlidingWindow) cleanupOldRequestsLoop() {
	ticker := time.NewTicker(1 * time.Second) // run cleanup every second
	for range ticker.C {
		cutoff := (time.Now().UnixMilli() - slidingWindow.Duration*2); // cutoff of the provious window

		// Step 1: Quickly collect up to 2 expired timestamps
		var oldTimestamps []int64;
		slidingWindow.mutex.Lock();
		count := 0
		for index,ts := range slidingWindow.RequestHistory {
			if ts <= cutoff {
				oldTimestamps = append(oldTimestamps, ts);
				slidingWindow.RequestHistory = slidingWindow.RequestHistory[index:];
				count++;
				if count >= 2 { // clean only 2 entries at a time
					break;
				}
			}
		}
		slidingWindow.mutex.Unlock();
		if len(oldTimestamps) > 0 {
			fmt.Println("Cleaned up", len(oldTimestamps), "old entries");
		}
	}
}

func(slidingWindow *SlidingWindow) addRequestHistory(now int64) {
	slidingWindow.mutex.Lock();
	slidingWindow.RequestHistory = append(slidingWindow.RequestHistory, now)
	slidingWindow.mutex.Unlock();
}

func(slidingWindow *SlidingWindow) IsRequestAllowed() bool {
	now := time.Now().UnixMilli();
	fmt.Println(fmt.Sprintf("The request time %d", now));
	if len(slidingWindow.RequestHistory) == 0 {
		slidingWindow.addRequestHistory(now);
		return true;
	}
	timeFrame:=now-slidingWindow.Duration;
	var reqCount int32 = 0;
	slidingWindow.mutex.Lock();
	defer slidingWindow.mutex.Unlock();
	fmt.Println(fmt.Sprintf("The last %d duration window start from %d", slidingWindow.Duration, timeFrame));
	for i := len(slidingWindow.RequestHistory) - 1; i >= 0; i-- {
		data := slidingWindow.RequestHistory[i]
		if data >= timeFrame {
			fmt.Println(fmt.Sprintf("The last request time %d in the %d timeFrame", data, timeFrame));
			reqCount++;
		}
		if reqCount>=slidingWindow.MaxCount {
			fmt.Println(fmt.Sprintf("The request count %d in last %d seconds", reqCount, slidingWindow.Duration/1000));
			return false;
		}
	}
	fmt.Println(fmt.Sprintf("The request count %d in last %d seconds", reqCount, slidingWindow.Duration/1000));
	if reqCount<slidingWindow.MaxCount {
		slidingWindow.addRequestHistory(now);
		return true;
	}
	return false;
}

func InitialiseSlidingWindow() SlidingWindow {
	slidingWindow := SlidingWindow{MaxCount: 2, Duration: 10000, RequestHistory: []int64{}};
	return slidingWindow;
}