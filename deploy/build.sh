#!/bin/sh

set -xe

APP=$1

cd $(dirname $0)/..
ROOT=$(pwd)

cd cmd/$APP
go build

mkdir -p /opt/iot/
install -m 0755 $APP /opt/iot/$APP