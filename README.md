# Hochwassernachrichtendienst Exporter

an unofficial prometheus exporter for the *Hochwassernachrichtendienst Bayern*.

## Usage

```
Usage of ./hochwassernachrichtendienst_exporter:
      --listenAddress string   Address on which to expose metrics (default ":8777")
      --metricsPath string     Path under which to expose metrics (default "/metrics")
```

## Prometheus Configuration

The exporter fetches the station on-demand and therefor needs it the `stationId` as a HTTP parameter.
Sample curl call: `curl http://127.0.0.1:8777/metrics\?station\=10092000 -v`

The `station` parameter can be set by using relabeling. You can find the numeric station id in the respective
URL of the station page at https://www.hnd.bayern.de/pegel/meldestufen

### Example Config

```
scrape_configs:
  - job_name: 'hochwassernachrichtendienst'
    metrics_path: /metrics
    static_configs:
      - targets:
        - 10092000
        - 16005701
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_station
      - source_labels: [__param_station]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:8777  # The address:port on which hochwassernachrichtendienst_exporter can be reached
```

## Exposed values

Example output for station **16005701**.

```
# HELP current_warning_level current reported warning level
# TYPE current_warning_level gauge
current_warning_level{id="16005701",name="München Isar"} 0
# HELP hundered_year_flood_level_centimeters an event that reaches or surpasses that level with a probability of 1% per year, -1 if not present
# TYPE hundered_year_flood_level_centimeters gauge
hundered_year_flood_level_centimeters{id="16005701",name="München Isar"} 510
# HELP last_level_centimeters last level in centimeters
# TYPE last_level_centimeters gauge
last_level_centimeters{id="16005701",name="München Isar"} 130
# HELP last_outflow_cubicmeters_per_second last outflow in cubicmeters/second, -1 if not present
# TYPE last_outflow_cubicmeters_per_second gauge
last_outflow_cubicmeters_per_second{id="16005701",name="München Isar"} 88.9
# HELP warning_level_1_centimeters warning level 1 in centimeters, -1 if not present
# TYPE warning_level_1_centimeters gauge
warning_level_1_centimeters{id="16005701",name="München Isar"} 240
# HELP warning_level_2_centimeters warning level 2 in centimeters, -1 if not present
# TYPE warning_level_2_centimeters gauge
warning_level_2_centimeters{id="16005701",name="München Isar"} 300
# HELP warning_level_3_centimeters warning level 3 in centimeters, -1 if not present
# TYPE warning_level_3_centimeters gauge
warning_level_3_centimeters{id="16005701",name="München Isar"} 380
# HELP warning_level_4_centimeters warning level 4 in centimeters, -1 if not present
# TYPE warning_level_4_centimeters gauge
warning_level_4_centimeters{id="16005701",name="München Isar"} 520
* Connection #0 to host 127.0.0.1 left intact
```

## Notice

Do not use this software to protect life or property. The parsed HTML output may not be stable. Alerts might stop working at any moment.

## Copyright & License

2021 Maximilian Güntner <code@mguentner.de>

MIT
