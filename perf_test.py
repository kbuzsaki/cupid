from datetime import datetime
from subprocess import Popen, PIPE, DEVNULL
import statistics
import sys
import time

CLIENT = "perf-client"

ADDRESSES = [
    "127.0.0.1:12380",
    "127.0.0.1:22380",
    "127.0.0.1:32380",
]

def launch_publisher_proc(addresses, topic, count, publishers, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-publish", "-topic", topic, "-count", str(count), "-pubs", str(publishers)]
    daemon = Popen(args, stdout=PIPE)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_subscriber_proc(addresses, topic, subscribers, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-topic", topic, "-subs", str(subscribers)]
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_daemons(launcher, count, *args, **kwargs):
    daemons = []
    for _ in range(count):
        daemons.append(launcher(*args, **kwargs))
    return daemons

def do_simple_perf_test(publishers, subscribers, topic, messages, sprocs, pprocs):
    subscribe_daemons = launch_daemons(launch_subscriber_proc, sprocs, ADDRESSES, topic, subscribers)

    start = datetime.now()
    publish_daemons = launch_daemons(launch_publisher_proc, pprocs, ADDRESSES, topic, messages, publishers)

    for pd in publish_daemons:
        pd.wait()
    end = datetime.now()

    delta = end - start
    print("time:", delta)

def do_keep_alive_perf_test(subscribers, topic, procs):
    subscribe_daemons = launch_daemons(launch_subscriber_proc, procs, ADDRESSES, topic, subscribers, stdout=PIPE)
    time.sleep(30)
    pd = launch_daemons(launch_publisher_proc, 1, ADDRESSES, topic, 0, 1)
    pd[0].wait()

    all_measurements = []
    for sd in subscribe_daemons:
        stdout, stderr = sd.communicate()
        output = stdout.decode("utf8")
        print(output)
        measurements = output.split('\n')[:-2] #last measurement is from the publisher
        print(measurements)
        measurements = [removeSuffix(x) for x in measurements] #drop the 's' from the time
        all_measurements.extend(measurements)

    print_stats(all_measurements)

def launch_locker(count, addresses, topic, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-locker", "-addrs", addrstr, "-topic", topic, "-count", str(count)]
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init locker")
    return daemon

#takes in a time and strips the suffix and converts to ms if in s
def removeSuffix(duration):
    if "ms" in duration: 
        return float(duration[:-2])

    #suffix is 's'
    duration = duration[:-1]
    duration = float(duration)
    duration *= 1000
    return duration

def do_lock_perf_test(count, topic, procs):
    locker_daemons = launch_daemons(launch_locker, procs, count, ADDRESSES, topic, stdout=PIPE)
    all_measurements = []
    for ld in locker_daemons:
        stdout, stderr = ld.communicate()
        measurements = stdout.decode("utf8").split("\n")[:-1]
        all_measurements.extend([removeSuffix(x) for x in measurements])
    print_stats(all_measurements)

def print_stats(measurements):
    print("mean:", statistics.mean(measurements))
    print("stdev:", statistics.stdev(measurements))
    print("median:", statistics.median(measurements))
    print("max:", max(measurements))
    print("min:", min(measurements))

def maybe_help(target_len, actual_len):
    should_print = target_len is not actual_len
    if should_print:
        print_help()
    return should_print

def print_help():
    print("python3 perf_test.py -locker [topic] [iterations] [procs]")
    print("python3 perf_test.py -pubsub [topic] [messages] [subscribers] [publishers] [sprocs] [pprocs]")
    print("python3 perf_test.py -keepalive [topic] [iterations] [procs]")

if __name__ == "__main__":

    if len(sys.argv) is 1:
        print_help()
        exit()

    cmd = sys.argv[1]

    if cmd == "-help":
        print_help()
        exit()

    if cmd == "-locker":
        if maybe_help(5, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        count = int(sys.argv[3])
        procs = int(sys.argv[4])
        do_lock_perf_test(count, topic, procs)
        exit()

    if cmd == "-pubsub":
        if maybe_help(8, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        messages = int(sys.argv[3])
        subscribers = int(sys.argv[4])
        publishers = int(sys.argv[5])
        sprocs = int(sys.argv[6])
        pprocs = int(sys.argv[7])
        do_simple_perf_test(publishers, subscribers, topic, messages, sprocs, pprocs)
        exit()

    if cmd == "-keepalive":
        if maybe_help(5, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        iterations = int(sys.argv[3])
        procs = int(sys.argv[4])
        do_keep_alive_perf_test(iterations, topic, procs)
        exit()
