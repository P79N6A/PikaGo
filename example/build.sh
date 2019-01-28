#!/usr/bin/env bash
set -x
RUN_NAME=Pika.Server

log()
{
    case $1 in
        DEBUG)
        echo $2
        ;;

        INFO)
        echo -e "\033[32m$2 \033[0m"
        ;;

        WARN)
        echo -e "\033[33m$2 \033[0m"
        ;;

        ERROR)
        echo -e "\033[31m$2 \033[0m"
        ;;

        *)
        ;;
    esac
}

rm -r output
mkdir -p output/bin output/log output/conf
cp -r conf output

log "DEBUG" "Init docker consul"
docker-compose -f docker-compose.yml up -d
if [ $? -eq 0 ]; then
    log "INFO" "Start docker consul successfully !"
else
    log "ERROR" "Start docker consul unsuccessfully !"
fi

go build -o output/bin/$RUN_NAME