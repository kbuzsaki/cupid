from datetime import datetime
from subprocess import Popen, PIPE, DEVNULL
import sys
import time

CLIENT = "perf-client"

ADDRESSES = [
    "127.0.0.1:12380",
    "127.0.0.1:22380",
    "127.0.0.1:32380",
]

def launch_publisher(addresses, topic, count):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-publish", "-topic", topic, "-count", str(count)]
    daemon = Popen(args)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_subscriber(addresses, topic):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-topic", topic]
    daemon = Popen(args, stdout=DEVNULL, stderr=DEVNULL)
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

if __name__ == "__main__":
    topic = sys.argv[1]
    count = int(sys.argv[2])
    do_simple_perf_test(1, 200, topic, count)
