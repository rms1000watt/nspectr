#!/usr/bin/env bash

set -e

echo "Building go binary"
go build

echo "Building docker"
docker build -t rms1000watt/nspectr .