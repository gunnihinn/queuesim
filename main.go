package main

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Copyright 2019, Gunnar Þór Magnússon

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	exprand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func help(def Defaults) string {
	return strings.TrimSpace(fmt.Sprintf(`
queuesim - Simulate a simple bounded queue

USE:

    queuesim [OPTION]...

queuesim simulates a simple bounded queue using discrete ticks for time. It
simulates a single producer and consumer of the queue, where waiting requests
can time out, and keeps track of the successes, timeouts, and rejections.

It can be helpful to imagine a single tick being one millisecond long when
setting values for the various program options.

OPTIONS:

    -method=METHOD      Method to use when popping items from queue (default: %s)
                        Accepted values: FIFO, FILO, RANDOM.
    -rate=RATE          A new request comes every RATE ticks (default %d)
    -size=SIZE          Size of queue (default %d)
    -timeout=TIMEOUT    Requests have TIMEOUT ticks to complete (default %d)
    -work=WORK          Requests take WORK ticks to complete (default %d)
    -ticks=TICKS        Run for TICKS ticks (default %d)
    -h, -help           Print help and exit
    -version            Print version and exit

CONTRIBUTING:

Patches are welcome on the project's GitHub page:

	https://www.github.com/gunnihinn/queuesim

COPYRIGHT:

This software is licensed under the GPLv3.
Copyright 2019, Gunnar Þór Magnússon <gunnar@magnusson.io>.
`, def.method, def.rate, def.size, def.timeout, def.work, def.ticks))
}

const VERSION = "2"

const (
	// Methods to use when popping items from a queue.
	FIFO = iota
	FILO
	RAND
)

// Simulation models a run of a queue.
type Simulation struct {
	queue   []*Request
	size    int
	rate    int
	work    func() int
	timeout int
	method  int

	counter Result
	current *Request
	tick    int
}

// Result holds the result of simulating a queue run.
type Result struct {
	successes  int
	timeouts   int
	rejections int
}

// Availability computes the availability of the queue over the run.
func (r Result) Availability() float64 {
	return float64(r.successes) / float64(r.successes+r.timeouts+r.rejections)
}

// Config holds the parameters of a queue simulation.
type Config struct {
	size    int
	rate    int
	work    func() int
	timeout int
	method  int
}

// NewSimulation creates a new simulation from a configuration.
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

// Run runs a queue simulation.
func (s *Simulation) Run(ticks int) Result {
	for t := 0; t < ticks; t++ {
		s.Tick()
	}

	return s.counter
}

// Tick advances a simulation by a single tick.
func (s *Simulation) Tick() {
	s.tick++

	// Work on current request
	if s.current != nil {
		s.current.Tick()
		s.current.Work()

		if s.current.Timedout() {
			s.counter.timeouts++
			s.current = nil
		} else if s.current.Done() {
			s.counter.successes++
			s.current = nil
		}
	}

	// Time out queued requests
	nq := make([]*Request, 0, s.size)
	for _, r := range s.queue {
		r.Tick()
		if r.Timedout() {
			s.counter.timeouts++
		} else {
			nq = append(nq, r)
		}
	}
	s.queue = nq

	// Accept or reject incoming request
	if s.tick%s.rate == 0 {
		if len(s.queue) == s.size {
			s.counter.rejections++
		} else {
			s.queue = append(s.queue, NewRequest(s.work(), s.timeout))
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

// Request is an item in the queue.
type Request struct {
	work   int
	budget int
}

// NewRequest creates a new Request.
func NewRequest(work int, timeout int) *Request {
	return &Request{work, timeout}
}

// Tick advances a request by a single tick.
func (r *Request) Tick() {
	r.budget--
}

// Work acheives one unit of work on a Request.
func (r *Request) Work() {
	r.work--
}

// Timedout reports whether a Request has timed out.
func (r Request) Timedout() bool { return r.budget <= 0 }

// Done reports whether a Request is done.
func (r Request) Done() bool { return r.work <= 0 }

// Defaults holds the main program's default parameter values.
type Defaults struct {
	rate    int
	timeout int
	work    int
	size    int
	method  string
	ticks   int
}

func main() {
	defaults := Defaults{
		rate:    10,
		timeout: 100,
		work:    30,
		size:    5,
		method:  "FIFO",
		ticks:   100000,
	}

	flags := struct {
		rate    *int
		timeout *int
		work    *int
		size    *int
		method  *string
		ticks   *int
		raw     *bool
		help    *bool
		h       *bool
		version *bool
	}{
		flag.Int("rate", defaults.rate, "A new request comes every RATE ticks"),
		flag.Int("timeout", defaults.timeout, "Requests have TIMEOUT ticks to complete"),
		flag.Int("work", defaults.work, "Requests take WORK ticks to complete"),
		flag.Int("size", defaults.size, "Size of queue"),
		flag.String("method", defaults.method, "Method to use when popping elements from queue"),
		flag.Int("ticks", defaults.ticks, "Number of ticks to run for"),
		flag.Bool("raw", false, "Don't pretty-print result"),
		flag.Bool("help", false, "Print help and exit"),
		flag.Bool("h", false, "Print help and exit"),
		flag.Bool("version", false, "Print version and exit"),
	}
	flag.Parse()

	if *flags.h || *flags.help {
		fmt.Println(help(defaults))
		return
	}

	if *flags.version {
		fmt.Println(VERSION)
		return
	}

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

	w := distuv.Poisson{
		Lambda: 1,
		Src:    exprand.NewSource(uint64(time.Now().Unix())),
	}

	work := func() int {
		return int(float64(*flags.work) * math.Pow(2, w.Rand()))
	}

	sim := NewSimulation(Config{
		rate:    *flags.rate,
		timeout: *flags.timeout,
		work:    work,
		size:    *flags.size,
		method:  method,
	})

	res := sim.Run(*flags.ticks)

	if *flags.raw {
		fmt.Printf("%f\n", res.Availability())
	} else {
		fmt.Printf("Availability: %.02f%%\n", 100*res.Availability())
	}
}
