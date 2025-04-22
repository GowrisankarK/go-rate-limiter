package algorithm

import "fmt"
import "time"

type SlidingWindow struct {
	MaxCount int32
	Duration int64
	LastRequestTimestamp int64
	RequestHistory []int64 // milliseconds
}

func(slidingWindow *SlidingWindow) cleanupOldRequests(timeStamp int64) {
	fmt.Println("Clean up started");
	for index,data:= range(slidingWindow.RequestHistory) {
		if(data>=timeStamp) {
			slidingWindow.RequestHistory = slidingWindow.RequestHistory[index:];
			return;
		}
	}
}

func(slidingWindow *SlidingWindow) addRequestHistory(now int64) {
	slidingWindow.LastRequestTimestamp = now
	slidingWindow.RequestHistory = append(slidingWindow.RequestHistory, now)
}

func(slidingWindow *SlidingWindow) IsRequestAllowed() bool {
	now := time.Now().UnixMilli();
	fmt.Println(fmt.Sprintf("The request time %d", now));
	if slidingWindow.LastRequestTimestamp == 0 {
		slidingWindow.addRequestHistory(now);
		return true;
	}
	timeFrame:=now-slidingWindow.Duration;
	defer slidingWindow.cleanupOldRequests(timeFrame);
	var reqCount int32 = 0;
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
	slidingWindow := SlidingWindow{MaxCount: 2, Duration: 10000, LastRequestTimestamp: 0, RequestHistory: []int64{}};
	return slidingWindow;
}