# pm_exporter

[![Build](https://github.com/GPORTALcloud/pm_exporter/actions/workflows/pm_exporter.yml/badge.svg?branch=main)](https://github.com/GPORTALcloud/pm_exporter/actions/workflows/pm_exporter.yml)

!!! STILL UNDER DEVELOPMENT !!!

pm_exporter is providing metrics fetched from platform management endpoints (such as iDRAC or ILO) and exposing those through HTTP for making those readable for prometheus.

For fetching those details the platform management specific API's are used (such as Redfish for iDRAC).
Those metrics are getting persisted within the runtime for the duration specified through command line arguments.

The list of platform managements to fetch details from is stored within a yaml file, each platform management needs to be defined there,
including their credentials, address and types. In addition, this config file allows adding
custom labels to the metrics which can be used for example for passing the id of the server within your
own inventory system or the switch the platform management is connected to.

Metrics which are getting collected:

| metric                         | Type    | description                                     |
| ------------------------------ |:-------:| :-----------------------------------------------|
| pm_platform_management_up      | gauge   | Defines if the platform management is reachable |
| pm_powersupply_health          | gauge   | Indicates the health status of the Power Supply |
| pm_battery_health              | gauge   | Indicates the health status of the Battery      |
| pm_cpu_health                  | gauge   | Indicates the health status of the CPU          |
| pm_fan_health                  | gauge   | Indicates the health status of the Fan          |
| pm_storage_health              | gauge   | Indicates the health status of the Storage      |
| pm_temperature_health          | gauge   | Indicates the health status of the Temperature  |
| pm_intrusion_health            | gauge   | Indicates if the chassis is closed or not       |
| pm_license_health              | gauge   | Indicates if the license is still valid         |
| pm_memory_health               | gauge   | Indicates the health status of the Memory       |
| pm_overall_health              | gauge   | Indicates if any of the health metrics are bad  |

Platform Managements are getting implemented once there is a host needs to get monitored, this list will grow with time.

Currently supported:
* iDRAC (redfish)


## Install
There are different ways of running the pm_exporter.

### Precompiled binaries
Precompiled binaries are getting created with each workflow run. Just download the artifact below "Actions" and
place them wherever you want.

### Docker images
pm_exporter is getting build and published to the Docker and Github registry.

You can just launch your own pm_exporter instance using the following command (arguments you're fine with the default you can just omit)

```
docker run -d --name pm_exporter --rm \
    -v /local/config.yml:/etc/pm_exporter.yml \
    msniveau/pm_exporter \
    --config.file=/etc/pm_exporter.yml \
    --web.listen-address=0.0.0.0:9096 \
    --web.enable-lifecycle \
    --metrics.persist_duration=5m \
    --metrics.refresh_interval=1m \
    --worker.count=5
```

## Command line arguments
| argument                  | default                | description                                            |
| ------------------------- |:----------------------:| :------------------------------------------------------|
| --config.file             | /etc/pm_exporter.yml   | Config file location                                   |
| -metric.persist_duration  | 1m30s                  | Duration metrics are getting persisted                 |
| -metric.refresh_interval  | 1m                     | Interval metrics are getting fetched                   |
| -web.enable-lifecycle     | false                  | Adds the /-/reload endpoint for config reload          |
| -web.listen-address       | 0.0.0.0:9096           | Listen address the HTTP server runs on                 |
| -worker.count             | 10                     | The amount of workers used for fetching the metrics    |