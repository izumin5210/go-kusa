#!/bin/sh

set -eu

cp -v /usr/share/zoneinfo/$TZDATA /etc/localtime
echo -e "$RUN_AT\t$WORKDIR/kusa" >> /var/spool/cron/crontabs/root

crond -l 2 -f
