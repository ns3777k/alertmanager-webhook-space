# Alertmanager webhook for JetBrains Space

## Getting started

Before you begin, you need multiple things to have:

1. Have an account on [Space](https://www.jetbrains.com/space/)
2. Channel to send alerts
3. Application with `Enable client credentials flow` on (Administration -> Applications).

## Building from source

```shell
$ CGO_ENABLED=0 go build -o alertmanager-webhook-jetbrains-space ./cmd/alertmanager-webhook-space/main.go
```

## Running docker

## Running without docker

## Example alertmanager configuration
