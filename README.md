### Overview

This plugin provides a zulip connection handler to send messages when alerts are raised

### Files
 * bin/sensu-zulip-handler

## Usage example

### Help

**sensu-zulip-handler**
```
The Sensu Go Zulip handler for notifying a channel

Usage:
  sensu-zulip-handler [flags]

Flags:
  -m, --bot-email string   The bot that messages will be sent as
  -k, --bot-key string     The bot key
  -c, --channel string     The channel to post messages to (default "#general")
  -h, --help               help for sensu-zulip-handler
  -u, --zulip-url string   The zulip url to send messages to, defaults to value of ZULIP_URL env variable
```
