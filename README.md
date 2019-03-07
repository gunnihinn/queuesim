# Queuesim

Queuesim is a bounded queue simulator. It shows how the success/failure rate of
a bounded queue with a single producer and consumer evoles with different
parameters. These parameters are:

* The queue size
* The rate of incoming requests
* The time it takes to process a request
* The request timeouts
* The method used to pop requests off the queue

The basic point of this program is to show that FIFO queues perform very poorly
in systems under load.

## Use

```
$ git clone https://www.github.com/gunnihinn/queuesim
$ cd queuesim
$ make
$ queuesim [OPTION...]
```
The options `queuesim` knows about are:
```
-method             Method to use when popping items from queue
-rate=RATE          A new request comes every RATE ticks
-size=SIZE          Size of queue
-timeout=TIMEOUT    Requests have TIMEOUT ticks to complete
-work=WORK          Requests take WORK ticks to complete
-ticks=TICKS        Run for TICKS ticks
-h, -help           Print help and exit
-version            Print version and exit
```
Acceptable `method` values are `FIFO`, `FILO` and `RANDOM`.

The program uses discrete ticks for the passage of time. It can be helpful to
imagine a single tick being one millisecond long when setting values for the
various program options.

## Try

Compare these runtime values:
```
$ queuesim -rate 10 -timeout 100 -work 30 -size 5 -method fifo
$ queuesim -rate 10 -timeout 100 -work 30 -size 5 -method filo
$ queuesim -rate 10 -timeout 100 -work 30 -size 5 -method random
```
That is, we receive a request every 10 ticks, to a queue that can hold 5 items.
Each request times out after 100 ticks, and it takes 30 ticks to process each
request.

You'll notice that the amount of successful to failed requests varies quite a
bit with the method we use to pop items from the queue. A FIFO queue performs
much worse than a FILO or random queue.
