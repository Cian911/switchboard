FROM alpine

COPY switchboard /usr/local/bin/switchboard

ENTRYPOINT ["switchboard"]
