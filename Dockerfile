FROM alpine:3.12.3
COPY alertmanager-webhook-space /
ENTRYPOINT ["/alertmanager-webhook-space"]
