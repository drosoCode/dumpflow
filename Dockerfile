FROM alpine:latest

RUN apk add --update ca-certificates p7zip
COPY dumpflow /dumpflow

ENTRYPOINT [ "/dumpflow" ]