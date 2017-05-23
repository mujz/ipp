FROM debian:jessie
MAINTAINER Mujtaba Al-Tameemi <mujtaba.altameemi@gmail.com>

RUN apt-get update && apt-get install -y ca-certificates

ENV PORT 80

ENTRYPOINT ["/ipp"]
ADD bin/ipp /ipp
