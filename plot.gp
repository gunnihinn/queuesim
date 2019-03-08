set terminal svg size 800, 640 font "Helvetica,16"
set output "queue.svg"
set datafile separator ","

set style line 1 \
    linecolor rgb '#0060ad' \
    linetype 1 linewidth 2

set style line 2 \
    linecolor rgb '#dd181f' \
    linetype 1 linewidth 2

set style line 3 \
    linecolor rgb '#000000' \
    linetype 1 linewidth 2

set key right bottom

set title "Queue of size 5, with timeout of 200 ticks, and work duration of 30 ticks"
set xlabel "Ticks between requests"
set ylabel "Availability %"

plot 'queue.csv' using 1:2 with lines linestyle 1 title 'FIFO',\
     'queue.csv' using 1:3 with lines linestyle 2 title 'FILO',\
     'queue.csv' using 1:4 with lines linestyle 3 title 'RANDOM'
