#!/bin/bash

set -xe

NAME=$1

cd $(dirname $0)

if [[ -f "$NAME.yaml" ]]; then
    tt -var "name=$NAME" -extra "$NAME.yaml" -tpl simple.service.tmpl -out /etc/systemd/system/iot-${NAME}.service
else
    tt -var "name=$NAME" -tpl simple.service.tmpl -out /etc/systemd/system/iot-${NAME}.service
fi

systemctl daemon-reload
systemctl enable iot-${NAME}
systemctl restart iot-${NAME}