# Alertmanager webhook for JetBrains Space

## Getting started

Before you begin, you need multiple things to have:

1. Have an account on [Space](https://www.jetbrains.com/space/)
2. Channel to send alerts
3. Application with `Enable client credentials flow` on (`Administration` -> `Applications`).

## Building from source

```shell
$ CGO_ENABLED=0 go build -o alertmanager-webhook-jetbrains-space ./cmd/alertmanager-webhook-space/main.go
```

## Running docker

## Running without docker

If you're on linux and have systemd, here is a sample configuration:

```
[Unit]
Description=Alertmanager webhook jetbrains space
After=network-online.target

[Service]
Type=simple
User=myuser
Group=myuser
ExecStart=/usr/local/bin/alertmanager-webhook-jetbrains-space
SyslogIdentifier=alertmanager-webhook-jetbrains-space
Restart=always

[Install]
WantedBy=multi-user.target
```

Don't forget to put the binary on `/usr/local/bin/alertmanager-webhook-jetbrains-space`.

## Sample alertmanager configuration

Do not copy-paste this blindly, that's just an example:

```yaml
global:
  resolve_timeout: 3m
templates:
  - '/etc/alertmanager/templates/*.tmpl'

receivers:
  - name: webhook
    webhook_configs:
      - url: http://127.0.0.1:9091/api/v1/webhook

route:
  group_by:
    - cluster
    - alertname
  group_interval: 5m
  group_wait: 30s
  receiver: webhook
  repeat_interval: 4h

```
