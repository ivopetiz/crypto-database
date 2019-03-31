#!/bin/bash

set -e

if [ ! -f "/.influxdb_configured" ]; then
    ./set_influxdb.sh
fi
exit 0