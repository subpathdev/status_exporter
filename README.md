# status exporter

## Output
| name | description |
| ---- | ----------- |
| status\_ip\_address{interface="\<name of interface\>", ipAddress=\"\<ip Addresses\>\"} | static to one and shows the ip addresses |
| status\_service{name="\<name of service\>"} | show the state of an service -1 not found / execution error, 0 active, 1 inactive |
