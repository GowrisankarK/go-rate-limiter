package main

import "fmt"
import "time"
import "github.com/GowrisankarK/go-rate-limiter/algorithm"

func ValidateFixedWindowRateLimiter() {
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

func main() {
	ValidateFixedWindowRateLimiter();
}