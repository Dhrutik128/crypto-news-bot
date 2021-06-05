package safe

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type routine struct {
	goroutine func(chan bool)
	stop      chan bool
}
type RoutineCtx struct {
	Goroutine func(ctx context.Context, routine RoutineCtx)
	Ctx       context.Context
	Out       chan map[string][]string
	In        chan map[string][]string
	Ring      *RingBuffer
}

func NewRoutineWithContext(routine func(ctx context.Context, routine RoutineCtx), ctx context.Context) RoutineCtx {
	in := make(chan map[string][]string)
	out := make(chan map[string][]string, 1)
	rb := NewRingBuffer(in, out)
	r := RoutineCtx{Ctx: ctx, Goroutine: routine, In: in, Out: out, Ring: rb}
	return r
}

// Go starts a recoverable Goroutine
func Go(goroutine func()) {
	GoWithRecover(goroutine, defaultRecoverGoroutine)
}

// GoTick starts a scheduled recoverable Goroutine with certain ticks. if executeFirst is true, Goroutine is started before ticking.
// todo -- use GoTick instead of Go(). This will prevent multiple routines hammering digicert api.
func GoTick(goroutine func(), tick time.Duration, executeFirst bool) {
	if executeFirst {
		Go(goroutine)
	}
	Go(func() {
		for range time.Tick(tick) {
			Go(goroutine)
		}
	})

}

// GoWithRecover starts a recoverable Goroutine using given customRecover() function
func GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}

func defaultRecoverGoroutine(err interface{}) {
	log.Errorf("Error in Go routine: %s", err)
}
