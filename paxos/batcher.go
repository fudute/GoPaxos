package paxos

import (
	"time"
)

const BatcherDefaultDuration = time.Millisecond
const BatcherDefaultLimit = 100

type Batcher struct {
	duration time.Duration
	limit    int
	Out      chan BatchRequest
	In       chan Request
}

type BatchRequest struct {
	Reqs []Request
	Done chan struct{}
}

var batcher *Batcher

func init() {
	batcher = &Batcher{
		duration: BatcherDefaultDuration,
		limit:    BatcherDefaultLimit,
		Out:      make(chan BatchRequest),
		In:       make(chan Request),
	}

	batcher.Run()
}

func (b *Batcher) Run() {
	go func() {
		batchReqs := BatchRequest{
			Reqs: make([]Request, 0, b.limit),
			Done: make(chan struct{}),
		}
		for {
			select {
			case req := <-b.In:
				batchReqs.Reqs = append(batchReqs.Reqs, req)
				if len(batchReqs.Reqs) >= b.limit {
					b.Out <- batchReqs
					<-batchReqs.Done
					batchReqs.Reqs = batchReqs.Reqs[:0]
				}
			case <-time.After(b.duration):
				if len(batchReqs.Reqs) != 0 {
					b.Out <- batchReqs
					<-batchReqs.Done
					batchReqs.Reqs = batchReqs.Reqs[:0]
				}
			}
		}
	}()
}

func GetBatcherInstance() *Batcher {
	return batcher
}

func SetBatcher(duration time.Duration, limit int) {
	batcher.duration = duration
	batcher.limit = limit
}
