# Extended prometheus metrics of a system
This is a tool, which can be used to export status metrics to [Prometheus](https://prometheus.io/).

## Command line parameters
| command | description |
| ------- | ----------- |
| help | print the help |
| listen | is the TCP network address |

## Metrics
| name | description |
| ---- | ----------- |
| status\_ip\_address{interface="\<name of interface\>", ipAddress=\"\<ip Addresses\>\"} | static to one and shows the ip addresses |
| status\_service{name="\<name of service\>"} | show the state of an service -1 not found / execution error, 0 active, 1 inactive |

## Configuration example prometheus
This show an example of the prometheus configuration of this exporter:
```
- job_name: status
    scrape_interval: 5s
    metrics_path: /
    static_configs:
    - targets: ['localhost:9998']
```
## Installation
To install this exporter you can use the following command:
```
go get github.com/subpathdev/status_exporter
```
