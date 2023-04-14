# Sonus Exporter for Prometheus

The Sonus Exporter is a Golang-based project that implements the multi-target exporter pattern and exposes metrics for a Ribbon SBC (Session Border Controller) in a format that can be scraped by Prometheus.

## Features

- Multi-target exporter pattern: Can scrape metrics from multiple SBC instances.
- Exposes metrics in Prometheus format for easy ingestion by Prometheus.
- Can be easily integrated into a Prometheus-based monitoring system.

## Requirements

- Golang 1.16 or higher
- Ribbon SBC


# Scaling

A single instance of `sonus_exporter` can be run for thousands of devices.

# Usage

## Installation

Binaries can be downloaded from the [Github
releases](https://github.com/ringsq/sonus_exporter/releases) page and need no special installation.

We also provide a sample [systemd unit file](examples/systemd/sonus_exporter.service).

## Running

Start `sonus_exporter` as a daemon or from CLI:

```sh
./sonus_exporter
```

Visit http://localhost:9116/sonus?module=if_mib&target=1.2.3.4 where `1.2.3.4` is the IP or
FQDN of the sonus device to get metrics from and `if_mib` is the default module, defined
in `sonus.yml`.

## Configuration

The default configuration file name is `sonus.yml` (currently unused).

The username/password used to connect to the SBCs is configured via the `SONUS_USER` and `SONUS_PASSWORD` environment variables.  The same username/password combination is used for all target SBCs.

## Prometheus Configuration

`target` can be passed as a parameter through relabelling.

Example config:
```YAML
scrape_configs:
  - job_name: 'sonus'
    static_configs:
      - targets:
        - 192.168.1.2  # sonus device.
        - sbc.local # sonus device.
    metrics_path: /sonus
    params:
      module: [if_mib]
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9700  # The sonus exporter's real hostname:port.
```

Similarly to [blackbox_exporter](https://github.com/prometheus/blackbox_exporter),
`sonus_exporter` is meant to run on a few central machines and can be thought of
like a "Prometheus proxy".

### TLS and basic authentication

The sonus Exporter supports TLS and basic authentication. This enables better
control of the various HTTP endpoints.

To use TLS and/or basic authentication, you need to pass a configuration file
using the `--web.config.file` parameter. The format of the file is described
[in the exporter-toolkit repository](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

Note that the TLS and basic authentication settings affect all HTTP endpoints:
/metrics for scraping, /sonus for scraping sonus devices, and the web UI.



## License

The Sonus Exporter is licensed under the [Apache License](https://opensource.org/licenses/apache-2-0).