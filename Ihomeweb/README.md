# Ihomeweb Service

This is the Ihomeweb service

Generated with

```
micro new zufang/Ihomeweb --namespace=go.micro --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.web.Ihomeweb
- Type: web
- Alias: Ihomeweb

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./Ihomeweb-web
```

Build a docker image
```
make docker
```