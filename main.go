package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	// Methods to use when popping items from a queue.
	FIFO = iota
	FILO
	RAND
)

var (
	success = 0
	failure = 0
)

type Request struct {
	done chan struct{}
	dead bool
}

// New creates a new request that times out after duration timeout.
func New(timeout time.Duration) *Request {
	r := &Request{
		done: make(chan struct{}),
	}

	go func(req *Request) {
		deadline := time.After(timeout)

		select {
		case <-deadline:
			req.dead = true
			failure++
		case <-req.done:
			success++
		}
	}(r)

	return r
}

// Work for duration d.
func (r *Request) Work(d time.Duration) {
	time.Sleep(d)
	close(r.done)
}

// Queue is a queue of requests.
type Queue struct {
	q    []*Request
	size int
	mut  *sync.Mutex
}

// NewQueue creates a new queue.
func NewQueue(size int) Queue {
	return Queue{
		q:    make([]*Request, 0, size),
		size: size,
		mut:  &sync.Mutex{},
	}
}

// Push a request onto the queue.
func (q *Queue) Push(r *Request) error {
	q.mut.Lock()
	defer q.mut.Unlock()

	if len(q.q) == q.size {
		return fmt.Errorf("No space in queue")
	}

	q.clean()
	q.q = append(q.q, r)

	return nil
}

// Pop a request from the queue.
func (q *Queue) Pop(method int) (*Request, error) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.clean()

	if len(q.q) == 0 {
		return nil, fmt.Errorf("No items in queue")
	}

	var r *Request
	switch method {
	case FIFO:
		r = q.q[0]
		q.q = q.q[1:len(q.q)]

	case FILO:
		r = q.q[len(q.q)-1]
		q.q = q.q[0 : len(q.q)-1]

	case RAND:
		k := rand.Intn(len(q.q))
		r = q.q[k]
		nq := q.q[0:k]
		nq = append(nq, q.q[k+1:len(q.q)]...)
		q.q = nq

	default:
		panic("Nonsense pop method used")
	}

	return r, nil
}

// Clean timed-out requests from the queue.
func (q *Queue) clean() {
	cq := make([]*Request, 0, q.size)
	for _, r := range q.q {
		if !r.dead {
			cq = append(cq, r)
		}
	}

	q.q = cq
}

func main() {
	flags := struct {
		rate    *int
		timeout *time.Duration
		work    *time.Duration
		size    *int
		method  *string
	}{
		flag.Int("rate", 1, "Number of incoming requests per second"),
		flag.Duration("timeout", time.Second, "Request timeout"),
		flag.Duration("work", time.Second, "Request work duration"),
		flag.Int("size", 10, "Size of queue"),
		flag.String("method", "FIFO", "Method to use when popping elements from queue"),
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())

	method := 0
	m := strings.ToLower(*flags.method)
	if strings.Contains(m, "fifo") {
		method = FIFO
	} else if strings.Contains(m, "filo") {
		method = FILO
	} else if strings.Contains(m, "rand") {
		method = RAND
	} else {
		panic("Unknown pop method given")
	}

	timeout := *flags.timeout
	work := *flags.work
	queue := NewQueue(*flags.size)

	go func() {
		tick := time.Tick(time.Second / time.Duration(*flags.rate))
		for range tick {
			r := New(timeout)
			if err := queue.Push(r); err != nil {
				failure++
				close(r.done)
			}
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			fmt.Printf("Success: %d\n", success)
			fmt.Printf("Failure: %d\n", failure)
			fmt.Printf("===\n")
		}
	}()

	for {
		req, err := queue.Pop(method)
		if err != nil {
			continue
		}

		req.Work(work)
		if !req.dead {
			success++
		}
	}
}
