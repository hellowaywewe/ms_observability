#!/usr/bin/python

from bcc import BPF
from bcc.utils import printb

# define BPF program
prog = """
int hello(void *ctx) {
    bpf_trace_printk("Hello, World!\\n");
    return 0;
}
"""

# load BPF program
b = BPF(text=prog)
b.attach_kprobe(event="blk_account_io_done", fn_name="hello")

print("The blk_account_io_done kernel function is called by MinSpore lenet training job.")
print("%-18s %-16s %-6s %s" % ("TIME(s)", "COMM", "PID", "MESSAGE"))

while 1:
    try:
        (task, pid, cpu, flags, ts, msg) = b.trace_fields()
    except ValueError:
        continue
    except KeyboardInterrupt:
        exit()
    if task == b'python' :
        printb(b"%-18.9f %-16s %-6d %s" % (ts, task, pid, msg))