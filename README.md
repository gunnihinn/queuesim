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
--method    The method to use when popping requests from the queue (default: FIFO)
--rate      Number of incoming requests per second (default: 1)
--size      The size of the queue (default: 10)
--timeout   Request timeout (default 1s)
--work      Request work duration (default 1s)
```
Acceptable `method` values are `FIFO`, `FILO` and `RANDOM`.

## Try

Compare these runtime values:
```
$ queuesim -rate 2 -work 750ms -method fifo
$ queuesim -rate 2 -work 750ms -method filo
$ queuesim -rate 2 -work 750ms -method random
```
That is, we receive 2 requests per second to a queue that can hold 10 items.
Each request times out after 1 second, and it takes 750ms to process each
request.

You'll notice that the amount of successful to failed requests varies quite a
bit with the method we use to pop items from the queue. A FIFO queue performs
much worse than a FILO or random queue.
