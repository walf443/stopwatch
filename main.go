// Example:
//
// 	import (
// 		"github.com/walf443/stopwatch"
// 	)
// func main() {
//	flag.Parse()
// 	stopwatch.Watch("init")
//	...
//	stopwatch.Watch("finish")
// }
//
//	you can run following
//	go run target.go --stopwatch
//
//	This library is not groutine safe.
//
package stopwatch

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	prevTime time.Time
	prevLock sync.RWMutex
)
var (
	resetTime time.Time
	resetLock sync.RWMutex
)

var logger *log.Logger
var enabled *bool

func init() {
	logger = log.New(os.Stdout, "stopwatch", log.Lshortfile)
	now := time.Now()

	resetLock.Lock()
	resetTime = now
	resetLock.Unlock()

	prevLock.Lock()
	prevTime = now
	prevLock.Unlock()

	enabled = flag.Bool("stopwatch", false, "if enable, work stopwatch")
}

func Reset(printfFormat string) {
	if *enabled {
		now := time.Now()

		resetLock.Lock()
		resetTime = now
		resetLock.Unlock()

		prevLock.Lock()
		prevTime = now
		prevLock.Unlock()

		output(now, printfFormat)
	}
}

func Watch(printfFormat string) {
	if *enabled {
		now := time.Now()
		prevLock.RLock()
		isZero := prevTime.IsZero()
		prevLock.RUnlock()
		if !isZero {
			output(now, printfFormat)
		}

		prevLock.Lock()
		prevTime = now
		prevLock.Unlock()
	}
}

func output(now time.Time, printfFormat string) {
	infoline := fmt.Sprintf(printfFormat)

	prevLock.Lock()
	dPrev := now.Sub(prevTime)
	prevLock.Unlock()

	resetLock.RLock()
	dReset := now.Sub(resetTime)
	resetLock.RUnlock()

	_, file, line, _ := runtime.Caller(2)

	output := fmt.Sprintf("[%14s][%14s] %s at %s:%d", dReset, dPrev, infoline, file, line)
	logger.Println(output)
}
