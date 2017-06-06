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

def launch_publisher(addresses, topic, count, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-publish", "-topic", topic, "-count", str(count)]
    daemon = Popen(args, stdout=stdout)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_subscriber(addresses, topic, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-topic", topic]
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_daemons(launcher, count, *args, **kwargs):
    daemons = []
    for _ in range(count):
        daemons.append(launcher(*args, **kwargs))
    return daemons

def do_simple_perf_test(publishers, subscribers, topic, count):
    subscribe_daemons = launch_daemons(launch_subscriber, subscribers, ADDRESSES, topic)

    start = datetime.now()
    publish_daemons = launch_daemons(launch_publisher, publishers, ADDRESSES, topic, count)
    for pd in publish_daemons:
        pd.wait()
    end = datetime.now()

    delta = end - start
    print("time:", delta)

def do_keep_alive_perf_test(subscribers, topic):
    subscribe_daemons = launch_daemons(launch_subscriber, subscribers, ADDRESSES, topic, stdout=PIPE)
    time.sleep(30)
    pd = launch_daemons(launch_publisher, 1, ADDRESSES, topic, 0)
    pd[0].wait()

    all_measurements = []
    for sd in subscribe_daemons:
        stdout, stderr = sd.communicate()
        output = stdout.decode("utf8")
        measurements = output.split('\n')[:-2] #last measurement is from the publisher
        measurements = [float(x[:-1]) for x in measurements] #drop the 's' from the time
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

def do_lock_perf_test(count, topic):
    locker_daemons = launch_daemons(launch_locker, 1, count, ADDRESSES, topic, stdout=PIPE)
    ld = locker_daemons[0]
    measurements = []
    stdout, stderr = ld.communicate()
    measurements = stdout.decode("utf8").split("\n")[:-1]
    measurements = [removeSuffix(x) for x in measurements]
    print_stats(measurements)

def print_stats(measurements):
    print("mean:", statistics.mean(measurements))
    print("stdev:", statistics.stdev(measurements))
    print("median:", statistics.median(measurements))
    print("max:", max(measurements))
    print("min:", min(measurements))

if __name__ == "__main__":
    if sys.argv[1] == "-locker":
        topic = sys.argv[2]
        count = int(sys.argv[3])
        do_lock_perf_test(count, topic)
    else:
        topic = sys.argv[1]
        count = int(sys.argv[2])
        do_simple_perf_test(1, 200, topic, count)
        do_keep_alive_perf_test(count, topic)
