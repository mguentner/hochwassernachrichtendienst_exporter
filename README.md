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
# HELP hnd_current_warning_level current reported warning level
# TYPE hnd_current_warning_level gauge
hnd_current_warning_level{id="16005701",name="München Isar"} 0
# HELP hnd_hundered_year_flood_level_millimeters an event that reaches or surpasses that level with a probability of 1% per year, -1 if not present
# TYPE hnd_hundered_year_flood_level_millimeters gauge
hnd_hundered_year_flood_level_millimeters{id="16005701",name="München Isar"} 5100
# HELP hnd_last_level_millimeters last level in millimeters
# TYPE hnd_last_level_millimeters gauge
hnd_last_level_millimeters{id="16005701",name="München Isar"} 1290
# HELP hnd_last_outflow_cubicmeters_per_second last outflow in cubicmeters/second, -1 if not present
# TYPE hnd_last_outflow_cubicmeters_per_second gauge
hnd_last_outflow_cubicmeters_per_second{id="16005701",name="München Isar"} 87.5
# HELP hnd_warning_level_1_millimeters warning level 1 in millimeters, -1 if not present
# TYPE hnd_warning_level_1_millimeters gauge
hnd_warning_level_1_millimeters{id="16005701",name="München Isar"} 2400
# HELP hnd_warning_level_2_millimeters warning level 2 in millimeters, -1 if not present
# TYPE hnd_warning_level_2_millimeters gauge
hnd_warning_level_2_millimeters{id="16005701",name="München Isar"} 3000
# HELP hnd_warning_level_3_millimeters warning level 3 in millimeters, -1 if not present
# TYPE hnd_warning_level_3_millimeters gauge
hnd_warning_level_3_millimeters{id="16005701",name="München Isar"} 3800
# HELP hnd_warning_level_4_millimeters warning level 4 in millimeters, -1 if not present
# TYPE hnd_warning_level_4_millimeters gauge
hnd_warning_level_4_millimeters{id="16005701",name="München Isar"} 5200
```

## Notice

Do not use this software to protect life or property. The parsed HTML output may not be stable. Alerts might stop working at any moment.

## Copyright & License

2021 Maximilian Güntner <code@mguentner.de>

MIT
