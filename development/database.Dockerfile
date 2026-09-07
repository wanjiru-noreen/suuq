FROM alpine:3.20

RUN apk add --no-cache sqlite
VOLUME ["/data"]

CMD ["sh", "-c", "mkdir -p /data && touch /data/suuq.db && exec tail -f /dev/null"]
