package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// Methods to use when popping items from a queue.
	FIFO = iota
	FILO
	RAND
)

type Simulation struct {
	queue   []*Request
	size    int
	rate    int
	work    int
	timeout int
	method  int

	// request status counters
	counter struct {
		success int
		timeout int
		reject  int
	}

	current *Request
	tick    int
}

type Config struct {
	size    int
	rate    int
	work    int
	timeout int
	method  int
}

func NewSimulation(cfg Config) *Simulation {
	switch cfg.method {
	case FIFO, FILO, RAND:

	default:
		panic("Unsupported pop method given")
	}

	return &Simulation{
		queue:   make([]*Request, 0, cfg.size),
		size:    cfg.size,
		rate:    cfg.rate,
		work:    cfg.work,
		timeout: cfg.timeout,
		method:  cfg.method,
	}
}

func (s *Simulation) Tick() {
	s.tick++

	// Work on current request
	if s.current != nil {
		s.current.Tick()
		s.current.Work()

		if s.current.Timedout() {
			s.counter.timeout++
		} else if s.current.Done() {
			s.counter.success++
			s.current = nil
		}
	}

	// Time out queued requests
	nq := make([]*Request, 0, s.size)
	for _, r := range s.queue {
		r.Tick()
		if r.Timedout() {
			s.counter.timeout++
		} else {
			nq = append(nq, r)
		}
	}
	s.queue = nq

	// Accept or reject incoming request
	if s.tick%s.rate == 0 {
		if len(s.queue) == s.size {
			s.counter.reject++
		} else {
			s.queue = append(s.queue, NewRequest(s.work, s.timeout))
		}
	}

	// Pick new request to work on
	if s.current == nil && len(s.queue) > 0 {
		switch s.method {
		case FIFO:
			s.current = s.queue[0]
			s.queue = s.queue[1:len(s.queue)]

		case FILO:
			s.current = s.queue[len(s.queue)-1]
			s.queue = s.queue[0 : len(s.queue)-1]

		case RAND:
			k := rand.Intn(len(s.queue))
			s.current = s.queue[k]
			nq := s.queue[0:k]
			nq = append(nq, s.queue[k+1:len(s.queue)]...)
			s.queue = nq

		default:
			panic("Nonsense pop method used")
		}
	}
}

type Request struct {
	work   int
	budget int
}

func NewRequest(work int, timeout int) *Request {
	return &Request{work, timeout}
}

func (r *Request) Tick() {
	r.budget--
}

func (r *Request) Work() {
	r.work--
}

func (r Request) Timedout() bool { return r.budget <= 0 }

func (r Request) Done() bool { return r.work <= 0 }

func main() {
	flags := struct {
		rate    *int
		timeout *int
		work    *int
		size    *int
		method  *string
		ticks   *int
	}{
		flag.Int("rate", 10, "A new request comes every RATE ticks"),
		flag.Int("timeout", 100, "Requests have TIMEOUT ticks to complete"),
		flag.Int("work", 30, "Requests take WORK ticks to complete"),
		flag.Int("size", 5, "Size of queue"),
		flag.String("method", "FIFO", "Method to use when popping elements from queue"),
		flag.Int("ticks", 100000, "Number of ticks to run for"),
	}
	flag.Parse()

	if *flags.rate <= 0 {
		panic("Rate must be positive")
	}

	if *flags.timeout <= 0 {
		panic("Timeout must be positive")
	}

	if *flags.work <= 0 {
		panic("Work must be positive")
	}

	if *flags.size <= 0 {
		panic("Size must be positive")
	}

	method := 0
	m := strings.ToLower(*flags.method)
	if strings.Contains(m, "fifo") {
		method = FIFO
	} else if strings.Contains(m, "filo") {
		method = FILO
	} else if strings.Contains(m, "rand") {
		rand.Seed(time.Now().Unix())
		method = RAND
	} else {
		panic("Unknown pop method given")
	}

	cfg := Config{
		rate:    *flags.rate,
		timeout: *flags.timeout,
		work:    *flags.work,
		size:    *flags.size,
		method:  method,
	}
	sim := NewSimulation(cfg)

	for t := 0; t < *flags.ticks; t++ {
		sim.Tick()
	}

	fmt.Printf("Success: %d\n", sim.counter.success)
	fmt.Printf("Timeout: %d\n", sim.counter.timeout)
	fmt.Printf("Reject: %d\n", sim.counter.reject)
}
