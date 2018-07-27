# Nspectr

## Introduction

Nspectr ("Inspector") is a reverse proxy that will block any payloads containing "XXXXXX"

This is just a PoC project

## Contents

- [Install](#install)
- [Build](#build)
- [Run](#run)
- [cURL](#curl)
- [References](#references)

## Install

```bash
# Govendor
go get -u github.com/kardianos/govendor
govendor sync
```

## Build

```bash
./build.sh
```

## Run

```bash
docker-compose up -d
```

## cURL

```bash
# Success
curl localhost:7100

# Success
curl -d 'XXX' localhost:7100

# Fail
curl -d 'XXXXXXXXXXXXXXXX' localhost:7100
```

## References

- https://gist.github.com/jbardin/821d08cb64c01c84b81a
- https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
