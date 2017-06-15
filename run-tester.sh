#!/bin/bash

for repl in 0 1 3 5; do
	for gr in 1 10; do
		for procs in 1 10 20 50 100 200; do
			fname=nop_repl_${repl}_gr_${gr}_procs_${procs} 
			./tester.sh 10 $repl $fname nop foo 0 $gr 10 $procs 
		done
	done
done

