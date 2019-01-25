#!/usr/bin/env bash
set -x
RUN_NAME=Pika.Server

CUR_DIR=$(cd $(dirname $0); pwd)
CUR_DATE=$(date +%Y-%m-%d)
export CONF_DIR=$CUR_DIR/conf

kill -9 `pgrep -f Pika.Server`
exec ./output/bin/$RUN_NAME
