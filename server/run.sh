#!/bin/bash

chmod +x ./server
exec ./server -listen-address ":${PORT}"

