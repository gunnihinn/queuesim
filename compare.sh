#!/bin/bash

set -euo pipefail

size=200
timeout=40
work=20

for rate in {200..1}; do
    fifo=$(./queuesim -method fifo -size "$size" -timeout "$timeout" -work-avg "$work_avg" -work-var "$work_var" -rate "$rate" -raw)
    filo=$(./queuesim -method filo -size "$size" -timeout "$timeout" -work-avg "$work_avg" -work-var "$work_var" -rate "$rate" -raw)
    rand=$(./queuesim -method rand -size "$size" -timeout "$timeout" -work-avg "$work_avg" -work-var "$work_var" -rate "$rate" -raw)

    n1=$(expr 100*$fifo | bc)
    n2=$(expr 100*$filo | bc)
    n3=$(expr 100*$rand | bc)

    printf "%d,%.0f,%.0f,%.0f\n" $rate $n1 $n2 $n3
done
