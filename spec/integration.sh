#!/bin/bash

# Job control
set -m

# Start KV server
pushd ./server
	go build
	chmod +x ./server
	./server &
popd

# Return value to export once
nc -l 8080 < ./spec/simple.html &

# Start exporting data to KV server
pushd ./exporter
	go build
	chmod +x ./exporter
	./exporter &
popd

# Fail by default
exit_code=1

# Wait for key to have exported value
tries=0
while [ $tries -lt 10 ]; do
	if curl -v http://localhost:8889/key | grep VALUE; then
		exit_code=0
		break
	fi
	((tries++))
	sleep 5
done

trap 'kill -9 $(jobs -p)' EXIT

exit $exit_code
