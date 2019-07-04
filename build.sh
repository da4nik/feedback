#!/bin/bash

go build -o feedback -ldflags "-linkmode external -extldflags -static"
docker build -t 172.16.1.60:5000/feedback:6 .
docker push 172.16.1.60:5000/feedback:6
