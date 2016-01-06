FROM ubuntu:15.10
MAINTAINER Chris Snell <chris.snell@revinate.com>
ADD kube-journald-filter /kube-journald-filter
RUN apt-get update && apt-get -y install libsystemd-dev nmap

CMD ./kube-journald-filter -alt-journal-base=/journal | /usr/bin/ncat $LOG_DEST 514
