#!/usr/bin/env bash
# Stops the process if something fails
set -xe
touch /var/app/current/go.bak
sudo rm /var/app/current/go.*

# get all of the dependencies needed
go get -d ./...

# create the application binary that eb uses
GOOS=linux GOARCH=amd64 go build -o bin/application -ldflags="-s -w"