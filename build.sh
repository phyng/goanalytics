#!/usr/bin/env bash

CGO_ENABLED=0 go build server.go
sudo docker build -t phyng/goanalytics:latest ./
rm server
