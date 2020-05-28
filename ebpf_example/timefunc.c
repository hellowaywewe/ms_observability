#include <uapi/linux/ptrace.h>
#include <linux/blkdev.h>

BPF_HASH(start, struct request *);
BPF_HISTOGRAM(latency_dist);

int kprobe_start(struct pt_regs *ctx, struct request *req)
{
    u64 tsp = bpf_ktime_get_ns();
    start.update(&req, &tsp);
    return 0;
}

int kprobe_end(struct pt_regs *ctx, struct request *req)
{
    u64 *tsp, delta;
    tsp = start.lookup(&req);

    if (!tsp) {
        return 0;
    }

    //时间差
    delta = bpf_ktime_get_ns() - *tsp;
    delta /= 1000;

    latency_dist.increment(bpf_log2l(delta));

    start.delete(&req);
    return 0;
}