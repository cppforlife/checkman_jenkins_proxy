#!/bin/bash

# Used by buildpack https://github.com/shageman/buildpack-binary

chmod +x ./server

exec ./server                                                           \
	-listen-address            ":${PORT}"                                 \
	-http-discoverer-endpoint  "http://peer-example.cfapps.io/_discover"  \
	-http-discoverer-every     5s
