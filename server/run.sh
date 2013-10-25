#!/bin/bash

# Used by buildpack https://github.com/shageman/buildpack-binary

chmod +x ./server

exec ./server -listen-address ":${PORT}"
