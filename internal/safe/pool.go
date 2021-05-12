package safe

import (
	"context"
	"sync"
	"time"
)

// Pool is a pool of go routines
type Pool struct {
	routines    []routine
	routinesCtx []RoutineCtx
	waitGroup   sync.WaitGroup
	lock        sync.Mutex
	baseCtx     context.Context
	baseCancel  context.CancelFunc
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewPool creates a Pool
func NewPool(parentCtx context.Context) *Pool {
	baseCtx, baseCancel := context.WithCancel(parentCtx)
	ctx, cancel := context.WithCancel(baseCtx)
	return &Pool{
		baseCtx:    baseCtx,
		baseCancel: baseCancel,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (p *Pool) DropRoutineByReference(ref string) bool {
	for i, r := range p.routinesCtx {
		rr := r.Ctx.Value("ref")
		if ref == rr {
			p.routinesCtx = append(p.routinesCtx[:i], p.routinesCtx[i+1:]...)
			return true
		}
	}
	return false
}
func (p *Pool) GetRoutineByReference(ref string) *RoutineCtx {
	for _, r := range p.routinesCtx {
		rr := r.Ctx.Value("ref")
		if ref == rr {
			return &r
		}
	}
	return nil
}

// Ctx returns main context
func (p *Pool) Ctx() context.Context {
	return p.baseCtx
}

// AddGoCtx adds a recoverable Goroutine with a context without starting it
func (p *Pool) AddGoCtx(goroutine RoutineCtx) {
	p.lock.Lock()
	p.routinesCtx = append(p.routinesCtx, goroutine)
	p.lock.Unlock()
}

// GoCtx starts a recoverable Goroutine with a context
func (p *Pool) GoCtx(goroutine RoutineCtx) {
	p.lock.Lock()

	p.routinesCtx = append(p.routinesCtx, goroutine)
	p.lock.Unlock()
	p.waitGroup.Add(1)
	Go(func() {
		defer p.waitGroup.Done()
		goroutine.Goroutine(goroutine.Ctx, goroutine)
		goroutine.Ring.end = time.Now()
	})
}

// addGo adds a recoverable Goroutine, and can be stopped with stop chan
func (p *Pool) addGo(goroutine func(stop chan bool)) {
	p.lock.Lock()
	newRoutine := routine{
		goroutine: goroutine,
		stop:      make(chan bool, 1),
	}
	p.routines = append(p.routines, newRoutine)
	p.lock.Unlock()
}
func (p *Pool) GoTick(goroutine func(stop chan bool), tick time.Duration, executeFirst bool) {
	p.lock.Lock()
	newRoutine := routine{
		goroutine: goroutine,
		stop:      make(chan bool, 1),
	}
	p.routines = append(p.routines, newRoutine)
	GoTick(func() {
		goroutine(newRoutine.stop)
	}, tick, executeFirst)
	p.lock.Unlock()
}

// Go starts a recoverable Goroutine, and can be stopped with stop chan
func (p *Pool) Go(goroutine func(stop chan bool)) {
	p.lock.Lock()
	newRoutine := routine{
		goroutine: goroutine,
		stop:      make(chan bool, 1),
	}
	p.routines = append(p.routines, newRoutine)
	p.waitGroup.Add(1)
	Go(func() {
		defer p.waitGroup.Done()
		goroutine(newRoutine.stop)
	})
	p.lock.Unlock()
}

// Stop stops all started routines, waiting for their termination
func (p *Pool) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cancel()
	for _, routine := range p.routines {
		routine.stop <- true
	}
	p.waitGroup.Wait()
	for _, routine := range p.routines {
		close(routine.stop)
	}
}

func (p *Pool) Count() int {
	if p != nil && p.routines != nil {
		return len(p.routines)
	}
	return 0
}

// Cleanup releases resources used by the pool, and should be called when the pool will no longer be used
func (p *Pool) Cleanup() {
	p.Stop()
	p.lock.Lock()
	defer p.lock.Unlock()
	p.baseCancel()
}

// Start starts all stopped routines
func (p *Pool) Start() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.ctx, p.cancel = context.WithCancel(p.baseCtx)
	for i := range p.routines {
		p.waitGroup.Add(1)
		p.routines[i].stop = make(chan bool, 1)
		Go(func() {
			defer p.waitGroup.Done()
			p.routines[i].goroutine(p.routines[i].stop)
		})
	}

	for _, routine := range p.routinesCtx {
		p.waitGroup.Add(1)
		Go(func() {
			defer p.waitGroup.Done()
			routine.Goroutine(p.ctx, routine)
		})
	}
}
