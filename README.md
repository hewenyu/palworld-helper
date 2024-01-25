# palworld-helper

## Overview

The Palworld RCON Helper is a standalone service that monitors player joins and leaves in Palworld servers using a custom protocol over TCP called rcon (remote console). 

It uses the rcon to fetch player data from the server and optionally, kick or ban players based on configurable whitelist logic.

This helper is primarily used for automation in managing servers, allowing server admins to more efficiently control and manage their server player base.


## configs

config.yaml setting

```yaml
# RCON settings
rcon_settings:
  endpoint:  "You Endpoint Here"
  password:  "You Password Here"

# interval for archiving backups, in seconds. 
archive_time_seconds: 3600

interval: 5s
timeout: 1s
```

## WhiteList

auto create .txt file