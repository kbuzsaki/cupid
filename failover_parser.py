from collections import defaultdict
import sys

def parse_time(timestr):
    *rest, time, off, tz  = timestr.split(" ")
    hr, mn, sec = time.split(":")
    return "{0:.1f}".format(float(sec))

def load_data(fn):
    lines = open(fn).read().split("\n")

    data = []
    for line in lines:
        if " : " not in line:
            continue
        timestr, statustr = map(str.strip, line.split(" : "))
        time = parse_time(timestr)
        status = "success" in statustr
        data.append((time, status))
    return data

def load_all_data(names):
    datas = [dict(load_data(name)) for name in names]
    keys = sorted(sum(map(list, (data.keys() for data in datas)), []), key=float)

    all_data = {}
    for k in keys:
        all_data[k] = tuple(data.get(k) for data in datas)
    return {k: v for k, v in all_data.items() if None not in v}

if __name__ == "__main__":
    names = sys.argv[1:4]
    data = load_all_data(names)
    time_statuses = [(k, data[k]) for k in sorted(data.keys(), key=float)]


    last_primary = None
    first_failover = None

    primary_up = True
    failover_up = False
    for time, statuses in time_statuses:
        # the primary hasn't died yet
        if primary_up:
            if not any(statuses):
                primary_up = False
            else:
                last_primary = time
        # the failover hasn't taken over yet
        elif not failover_up:
            if any(statuses):
                first_failover = time
                failover_up = True
        # else the failover has taken over, so we don't care

    print("last_primary:", last_primary)
    print("first_failover:", first_failover)
    delta = float(first_failover) - float(last_primary)
    deltastr = "{0:.1f}".format(delta)
    print("failover period:", deltastr)
