#!/usr/bin/python
from __future__ import print_function
from bcc import BPF
from time import sleep, strftime
import sys
from bcc.utils import printb


b = BPF(src_file="timefunc.c")
b.attach_kprobe(event="blk_account_io_start", fn_name="kprobe_start")
b.attach_kprobe(event="blk_account_io_done", fn_name="kprobe_end")

if len(sys.argv) != 3:
     print(
 """
 Trace block device I/O latency, and print the distribution graph (histogram).

 Usage: %s [interval] [count]
 interval - time interval (seconds)
 count - how many times to record

 Example: Record once every second, record 10 times in total
 $ %s 1 10
 """ % (sys.argv[0], sys.argv[0]))
     sys.exit(1)

interval = int(sys.argv[1])
count = int(sys.argv[2])
print("Tracing block device I/O... Hit Ctrl-C to end.")

exiting = 0 if interval else 1
latency_dist = b.get_table("latency_dist")
while (1):
    try:
        sleep(interval)
    except KeyboardInterrupt:
        exiting = 1

    print()
    print("%-8s\n" % strftime("%H:%M:%S"), end="")

    latency_dist.print_log2_hist("us")
    latency_dist.clear()

    count -= 1
    if exiting or count == 0:
        exit()