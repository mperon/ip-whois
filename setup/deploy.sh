#!/usr/bin/env bash

DEST_SERVER=mperon
DEST_USER=devel
DEST_FOLDER=/opt/golang/bin/
DEST_CHOWN=devel:www-data
TODAY=$(date +"%Y%m%d_%H%M%S")
BASEDIR=$(realpath "${BASH_SOURCE%/*}/../")

rsync -avh ./bin/ip-whois-linux-amd64 ${DEST_USER}@${DEST_SERVER}:${DEST_FOLDER}

# run code on remote server
ssh -t ${DEST_USER}@${DEST_SERVER} <<COMMAND_TEXT
    echo "Reloading service.."
    sudo /bin/systemctl stop ip-whois
    sudo /bin/systemctl start ip-whois
COMMAND_TEXT
