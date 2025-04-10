package algorithm

import "fmt"
import "time"

type FixedWindow struct {
	MaxCount int32
	Duration int64
	StartTimestamp int64
	CurrentCount int32

}

func(fixedWindow *FixedWindow) IsRequestAllowed() bool {
	now := time.Now().UnixMilli();
	fmt.Println(fmt.Sprintf("The current time %d and difference %d", now, now-fixedWindow.StartTimestamp));
	if now-fixedWindow.StartTimestamp<=fixedWindow.Duration && fixedWindow.CurrentCount < fixedWindow.MaxCount{
		fixedWindow.CurrentCount++;
		return true;
	}
	return false;
}

func InitialiseFixedWindow() FixedWindow{
	fixedWindow := FixedWindow{MaxCount: 100, Duration: 10000, StartTimestamp: time.Now().UnixMilli(), CurrentCount: 0};
	return fixedWindow;
}