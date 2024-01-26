#!/bin/bash
# This script is used to start the endpoint service
# create config file from env vars

# yaml example:
# # RCON settings
# rcon_settings:
#   endpoint: "endpointHere"
#   password: "adminPasswordHere"

# # interval for archiving backups, in seconds.
# archive_time_seconds: 3600

# interval: 5s
# timeout: 1s

echo "Creating config file from env vars"
echo "rcon_settings:" > /app/config.yaml
echo "  endpoint: \"$ENDPOINT\"" >> /app/config.yaml
echo "  password: \"$PASSWORD\"" >> /app/config.yaml
echo "archive_time_seconds: 3600" >> /app/config.yaml
echo "interval: 5s" >> /app/config.yaml
echo "timeout: 1s" >> /app/config.yaml

# start the endpoint service
echo "Starting endpoint service"

# start service /app/monitor in /app
cd /app && ./monitor
