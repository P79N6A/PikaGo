#!/usr/bin/env bash

CUR_DIR=$(cd $(dirname $0); pwd)
export CONF_DIR=$CUR_DIR/conf

go run serve.go