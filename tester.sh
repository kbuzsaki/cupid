#!/bin/bash

trials=$1
replicas=$2
outputfile=$3
cmd=$4
shift
shift
shift
shift
 
for i in $(seq 1 $trials); do 
	./launch.sh $replicas &
	pid=$!
	sleep 2
	python3 perf_test.py -$cmd $@ >> data/$outputfile
	#kill -2 $pid
	pkill cupid-server
	make fresh
done

