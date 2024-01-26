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
echo "rcon_settings:" > /etc/endpoint.yaml
echo "  endpoint: \"$ENDPOINT\"" >> /etc/endpoint.yaml
echo "  password: \"$PASSWORD\"" >> /etc/endpoint.yaml
echo "archive_time_seconds: 3600" >> /etc/endpoint.yaml
echo "interval: 5s" >> /etc/endpoint.yaml
echo "timeout: 1s" >> /etc/endpoint.yaml

# start the endpoint service
echo "Starting endpoint service"

# start service /app/monitor in /app
cd /app && ./monitor