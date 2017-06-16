from datetime import datetime
from subprocess import Popen, PIPE, DEVNULL
import statistics
import sys
import time
import numpy

CLIENT = "perf-client"

ADDRESSES = [
    "127.0.0.1:12380",
    "127.0.0.1:22380",
    "127.0.0.1:32380",
    "127.0.0.1:42380",
    "127.0.0.1:52380"
]

def launch_publisher_proc(addresses, topic, numGoRoutines, publishers, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-publish", "-topic", topic, "-numGoRoutines", str(numGoRoutines), "-pubs", str(publishers)]
    print(args)
    daemon = Popen(args, stdout=PIPE)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_subscriber_proc(addresses, topic, subscribers, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-addrs", addrstr, "-topic", topic, "-subs", str(subscribers)]
    print(args)
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init publisher")
    return daemon

def launch_daemons(launcher, numGoRoutines, *args, **kwargs):
    daemons = []
    for _ in range(numGoRoutines):
        daemons.append(launcher(*args, **kwargs))
    return daemons

#takes in a csv and the number of procs and extends/shrinks topics that size
def get_topics_list(topics_str, procs):
    topic_list = topics_str.split(",")
    topic_arg = []
    for i in range(procs):
        index = i % len(topic_list)
        topic = topic_list[index]
        topic_arg.append(topic)

    return topic_arg

def do_simple_perf_test(publishers, subscribers, topics_str, messages, sprocs, pprocs):


    # subscribe_daemons = launch_daemons(launch_subscriber_proc, sprocs, ADDRESSES, topic, subscribers)

    subscribe_daemons = []
    topics_list = get_topics_list(topics_str, sprocs)
    for topic in topics_list:
        sd = launch_daemons(launch_subscriber_proc, 1, ADDRESSES, topic, subscribers)[0]
        subscribe_daemons.append(sd)

    topics_list = get_topics_list(topics_str, pprocs)
    publish_daemons = []
    for topic in topics_list:
        pd = launch_daemons(launch_publisher_proc, 1, ADDRESSES, topic, messages, publishers)[0]
        publish_daemons.append(pd)

    # publish_daemons = launch_daemons(launch_publisher_proc, pprocs, ADDRESSES, topic, messages, publishers)

    for sd in subscribe_daemons:
        sd.wait()

    all_measurements = []
    for pd in publish_daemons:
        stdout, stderr = pd.communicate()
        output = stdout.decode("utf8").split("\n")[:-1]
        all_measurements.extend(output)

    all_measurements = [float(x) for x in all_measurements]
    print("mean   :", toMilli(statistics.mean(all_measurements)))
    print("median :", toMilli(statistics.median(all_measurements)))

def do_keep_alive_perf_test(subscribers, topic, procs):
    subscribe_daemons = launch_daemons(launch_subscriber_proc, procs, ADDRESSES, topic, subscribers, stdout=PIPE)
    time.sleep(30)
    pd = launch_daemons(launch_publisher_proc, 1, ADDRESSES, topic, 0, 1)
    pd[0].wait()

    all_measurements = []
    for sd in subscribe_daemons:
        stdout, stderr = sd.communicate()
        output = stdout.decode("utf8")
        measurements = output.split('\n')[:-2] #last measurement is from the publisher
        all_measurements.extend(measurements)

    print_stats(all_measurements)

def launch_locker(numGoRoutines, iterations, addresses, topic, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-locker", "-addrs", addrstr, "-topic", topic, "-numGoRoutines", str(numGoRoutines), "-iterations", str(iterations)]
    print(args)
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init locker")
    return daemon

def do_lock_perf_test(numGoRoutines, iterations, topic, procs):
    locker_daemons = launch_daemons(launch_locker, procs, numGoRoutines, iterations, ADDRESSES, topic, stdout=PIPE)
    all_measurements = []
    for ld in locker_daemons:
        stdout, stderr = ld.communicate()
        measurements = stdout.decode("utf8").split("\n")[:-1]
        all_measurements.extend(measurements)
    print_stats(all_measurements)

def do_nop_perf_test(topic, numOps, iterations, goroutines, procs):
    nop_daemons = launch_daemons(launch_noper, procs, ADDRESSES, topic, numOps, iterations, goroutines, stdout=PIPE)
    all_measurements = []
    for nd in nop_daemons:
        stdout, stderr = nd.communicate()
        measurements = stdout.decode("utf8").split("\n")[:-1]
        all_measurements.extend(measurements)
    print_stats(all_measurements)

def launch_noper(addresses, topic, numOps, iterations, goroutines, stdout=DEVNULL):
    addrstr = ",".join(addresses)
    args = [CLIENT, "-nop", "-addrs", addrstr, "-topic", topic, "-numOps", str(numOps), "-numGoRoutines", str(goroutines), "-iterations", str(iterations)]
    print(args)
    daemon = Popen(args, stdout=stdout, stderr=DEVNULL)
    if daemon.poll():
        raise Exception("failed to init noper")
    return daemon

def toMilli(ts):
    return ts/1000000 #divide by 10^-6 to go from nano to milli


def print_stats(measurements):
    measurements = [float(x) for x in measurements]
    print("mean   :", toMilli(statistics.mean(measurements)))
    print("stdev  :", toMilli(statistics.stdev(measurements)))
    print("median :", toMilli(statistics.median(measurements)))
    print("max    :", toMilli(max(measurements)))
    print("min    :", toMilli(min(measurements)))
    print("90th p :", toMilli(numpy.percentile(measurements, 90)))
    print("99th p :", toMilli(numpy.percentile(measurements, 99)))
    print("total  :", toMilli(sum(measurements)))

def maybe_help(target_len, actual_len):
    should_print = target_len is not actual_len
    if should_print:
        print_help()
    return should_print

def print_help():
    print("python3 perf_test.py -locker     [topic1,topic2,topic3...] [numGoRoutines] [iterations] [procs]")
    print("python3 perf_test.py -pubsub     [topic] [numMessages] [numSubscribers] [numPublishers] [sprocs] [pprocs]")
    print("python3 perf_test.py -keepalive  [topic] [numGoRoutines] [procs]")
    print("python3 perf_test.py -nop        [topic] [numOps] [goroutines] [iterations] [procs]")

if __name__ == "__main__":

    if len(sys.argv) is 1:
        print_help()
        exit()

    cmd = sys.argv[1]

    if cmd == "-help":
        print_help()
        exit()

    if cmd == "-locker":
        if maybe_help(6, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        numGoRoutines = int(sys.argv[3])
        iterations = int(sys.argv[4])
        procs = int(sys.argv[5])
        do_lock_perf_test(numGoRoutines, iterations, topic, procs)

    if cmd == "-pubsub":
        if maybe_help(8, len(sys.argv)):
            exit()
        topics = sys.argv[2]
        messages = int(sys.argv[3])
        subscribers = int(sys.argv[4])
        publishers = int(sys.argv[5])
        sprocs = int(sys.argv[6])
        pprocs = int(sys.argv[7])
        do_simple_perf_test(publishers, subscribers, topics, messages, sprocs, pprocs)

    if cmd == "-keepalive":
        if maybe_help(5, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        iterations = int(sys.argv[3])
        procs = int(sys.argv[4])
        do_keep_alive_perf_test(iterations, topic, procs)

    if cmd == "-nop":
        if maybe_help(7, len(sys.argv)):
            exit()
        topic = sys.argv[2]
        numOps = int(sys.argv[3])
        goroutines = int(sys.argv[4])
        iterations = int(sys.argv[5])
        procs = int(sys.argv[6])
        do_nop_perf_test(topic, numOps, iterations, goroutines, procs)
