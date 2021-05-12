package safe

import (
	"strconv"
	"time"
)

type RingBuffer struct {
	inputChannel  <-chan map[string][]string
	outputChannel chan map[string][]string
	Processing    bool
	start         time.Time
	end           time.Time
}

func NewRingBuffer(inputChannel <-chan map[string][]string, outputChannel chan map[string][]string) *RingBuffer {
	rb := &RingBuffer{inputChannel: inputChannel, outputChannel: outputChannel, Processing: false, start: time.Now()}
	Go(func() {
		rb.Run()
	})
	return rb
}
func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		select {
		case r.outputChannel <- v:
			SetStatus(v, r)
		default:
			<-r.outputChannel
			SetStatus(v, r)
			r.outputChannel <- v
		}
	}
	close(r.outputChannel)
}
func SetStatus(v map[string][]string, r *RingBuffer) {
	v["processing"] = []string{strconv.FormatBool(r.Processing)}
	v["started"] = []string{r.start.String()}
	if !r.end.IsZero() {
		v["finished"] = []string{r.end.String()}
	}
}
