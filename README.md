# Heartbeatoshaker

This is a command line utility that pings a url or urls at intervals. The need for this tool is as a result of keeping websites that I host on fee hosting
platforms that shutdown the service due to inactivity.

## Installation

- Using go tool `go install github.com/gray-adeyi/heartbeatoshaker@latest`.
  Note: Installing it like this makes it available as `heartbeatoshaker` instead of `hbs`
- Install prebuilt binaries from the [Releases]()

## Usages

- `hbs` (Runs heartbeatoshaker with the config file it finds)
- `hbs init` (This creates a heartbeat.yaml file in the current path that you can configure)
- `hbs -url https://github.com/gray-adeyi/heartbeatoshaker  -interval 5m` (Pings the url at five minutes intervals)
- `hbs -config my-heartbeat.yml` (Runs heartbeatoshaker with the custom config file)