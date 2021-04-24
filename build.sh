#!/bin/sh
# Stops the process if something fails
sudo rm /var/app/current/go.*

cd src || exit

# get all of the dependencies needed
go get -d ./...

# create the application binary that eb uses
GOOS=linux GOARCH=amd64 go build -o ../bin/application -ldflags="-s -w"
