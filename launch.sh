#! /bin/bash

ONE_CLUSTER='http://127.0.0.1:12379'
THREE_CLUSTER='http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379'
FIVE_CLUSTER='http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379,http://127.0.0.1:42379,http://127.0.0.1:52379'



pids=""

function killservers() {
    echo
    echo "exiting..."
    kill -9 $pids
    echo done
}

trap "killservers" SIGINT

CLUSTER_ARG=""
if [ "$1" = "1" ]; then
    CLUSTER_ARG=$ONE_CLUSTER
elif [ "$1" = "3" ]; then
    CLUSTER_ARG=$THREE_CLUSTER
elif [ "$1" = "5" ]; then
    CLUSTER_ARG=$FIVE_CLUSTER
else
    echo "invalid arg number of replicas: '$1'"
    exit 1
fi

for id in $(seq 1 $1); do
    echo "launching $id"
    cupid-server --id $id --cluster $CLUSTER_ARG --port "$id"2380 &
    pids="$pids $!"
done

wait
