#!/bin/bash
cat > config.yaml << EOF
# RCON settings
rcon_settings:
  endpoint: ${RCON_ENDPOINT}
  password: ${RCON_PASSWORD}

# interval for archiving backups, in seconds.
archive_time_seconds: 3600

interval: 5s
timeout: 1s
EOF