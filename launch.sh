#! /bin/bash

NO_CLUSTER='none'
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

replicas=$1
cluster_arg=""
if [ "$1" = "0" ]; then
    replicas=1
    cluster_arg=$NO_CLUSTER
elif [ "$1" = "1" ]; then
    cluster_arg=$ONE_CLUSTER
elif [ "$1" = "3" ]; then
    cluster_arg=$THREE_CLUSTER
elif [ "$1" = "5" ]; then
    cluster_arg=$FIVE_CLUSTER
else
    echo "invalid arg number of replicas: '$1'"
    exit 1
fi

for id in $(seq 1 $replicas); do
    echo "launching $id"
    cupid-server --id $id --cluster $cluster_arg --port "$id"2380 &
    pids="$pids $!"
done

wait
