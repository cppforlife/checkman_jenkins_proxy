[![Build Status](https://travis-ci.org/cppforlife/checkman_jenkins_proxy.png?branch=master)](https://travis-ci.org/cppforlife/checkman_jenkins_proxy)

## Server

Toy KV store (currently must run as a single instance).

For Checkman: used to serve data exported from Jenkins.

```
cd server/
go build
chmod +x ./server
./server -h
```

```
curl -X GET http://localhost:8889/key
404 Not Found

curl -X PUT http://localhost:8889/key -d 'something'

curl -X GET http://localhost:8889/key
something

# ...expires after 30s by default...
curl -X GET http://localhost:8889/key
404 Not Found
```


## Exporter

Periodically fetches content from given url
and uploads it to another given url.

For Checkman: run this on your Jenkins machine and
let it fetch and export build data from Jenkins api.

```
cd exporter/
go build
chmod +x ./exporter
./exporter -h
```

```
./exporter -fetcher-endpoint "http://local-jenkins/api/json?depth=2" -store-key something
```


## Misc

```
while true; do nc -l 8080 < spec/simple.html; done
```
