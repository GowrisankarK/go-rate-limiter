package algorithm

import "fmt"
import "time"
import "sync"

type LeakyBucket struct {
	MaxTokenCount int32
	requestQueue []int32
	OutFillRate int64 //millisecond
	ticker     *time.Ticker
	stopChannel   chan struct{}
	mutex         sync.Mutex
}

func (leakyBucket *LeakyBucket) IsRequestAllowed(requestId int32) bool {
	leakyBucket.mutex.Lock();
	defer leakyBucket.mutex.Unlock();
	if int32(len(leakyBucket.requestQueue)) < leakyBucket.MaxTokenCount {
		leakyBucket.requestQueue=append(leakyBucket.requestQueue, requestId);
		return true;
	}
	return false;
}

func (leakyBucket *LeakyBucket) processRequestFromQueue() {
	leakyBucket.ticker = time.NewTicker(time.Duration(leakyBucket.OutFillRate) * time.Millisecond);
	leakyBucket.stopChannel = make(chan struct{})
	go func() {
		for {
			select {
				case <-leakyBucket.ticker.C :
					leakyBucket.mutex.Lock();
					fmt.Println(fmt.Sprintf("Process the request %d at %d", leakyBucket.requestQueue[0], time.Now().UnixMilli()));
					leakyBucket.requestQueue=leakyBucket.requestQueue[1:];
					leakyBucket.mutex.Unlock();
				case <-leakyBucket.stopChannel :
					leakyBucket.ticker.Stop();
			}
		}
	}()
}

func (leakyBucket *LeakyBucket) CloseRequestProcessing() {
	close(leakyBucket.stopChannel)
}

func InitialiseLeakyBucket() *LeakyBucket {
	leakyBucket := &LeakyBucket{MaxTokenCount: 3, requestQueue: []int32{}, OutFillRate: 10000};
	//Queue
	leakyBucket.processRequestFromQueue();
	return leakyBucket;
}