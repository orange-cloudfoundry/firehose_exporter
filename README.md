# Cloud Foundry Firehose Exporter [![Build Status](https://travis-ci.org/bosh-prometheus/firehose_exporter.png)](https://travis-ci.org/bosh-prometheus/firehose_exporter)

A [Prometheus][prometheus] exporter proxy for [Cloud Foundry Firehose][firehose] metrics. Please refer to the [FAQ][faq]
for general questions about this exporter.

## Architecture overview

![](https://cdn.rawgit.com/bosh-prometheus/firehose_exporter/master/architecture/architecture.svg)

## Installation

### Binaries

Download the already existing [binaries](https://github.com/bosh-prometheus/firehose_exporter/releases) for your
platform:

```bash
$ ./firehose_exporter <flags>
```

### From source

Using the standard `go install` (you must have [Go][golang] already installed in your local machine):

```bash
$ go install github.com/bosh-prometheus/firehose_exporter
$ firehose_exporter <flags>
```

### Docker

To run the firehose exporter as a Docker container, run:

```bash
$ docker run -p 9186:9186 boshprometheus/firehose-exporter <flags>
```

### Cloud Foundry

The exporter can be deployed to an already existing [Cloud Foundry][cloudfoundry] environment:

```bash
$ git clone https://github.com/bosh-prometheus/firehose_exporter.git
$ cd firehose_exporter
```

Modify the included [application manifest file][manifest] to include your [Cloud Foundry Firehose][firehose] properties.
Then you can push the exporter to your Cloud Foundry environment:

```bash
$ cf push
```

### BOSH

This exporter can be deployed using the [Prometheus BOSH Release][prometheus-boshrelease].

## Usage

### Flags

| Flag / Environment Variable | Required | Default | Description |
| --------------------------- | -------- | ------- | ----------- |
| `retro_compat.disable`<br />`FIREHOSE_EXPORTER_RETRO_COMPAT_DISABLE` | No | `False` | Disable retro compatibility |
| `retro_compat.enable_delta`<br />`FIREHOSE_EXPORTER_RETRO_COMPAT_ENABLE_DELTA` | No | `False` | Enable retro compatibility delta in counter |
| `metrics.shard_id`<br />`FIREHOSE_EXPORTER_DOPPLER_SUBSCRIPTION_ID` | No | `prometheus` | Cloud Foundry Nozzle Subscription ID |
| `metrics.expiration`<br />`FIREHOSE_EXPORTER_DOPPLER_METRIC_EXPIRATION` | No | `10 minutes` | How long Cloud Foundry metrics received from the Firehose are valid |
| `metrics.batch_size`<br />`FIREHOSE_EXPORTER_METRICS_BATCH_SIZE` | No | `infinite buffer` | Batch size for nozzle envelop buffer |
| `metrics.node_index`<br />`FIREHOSE_EXPORTER_NODE_INDEX` | No | `0` | Node index to use |
| `metrics.timer_rollup_buffer_size`<br />`FIREHOSE_EXPORTER_TIMER_ROLLUP_BUFFER_SIZE` | No | `0` | The number of envelopes that will be allowed to be buffered while timer http metric aggregations are running |
| `filter.deployments`<br />`FIREHOSE_EXPORTER_FILTER_DEPLOYMENTS` | No | | Comma separated deployments to filter |
| `filter.events`<br />`FIREHOSE_EXPORTER_FILTER_EVENTS` | No | | Comma separated events to filter. If not set, all events will be enabled (`ContainerMetric`, `CounterEvent`, `HttpStartStop`, `ValueMetric`) |
| `logging.url`<br />`FIREHOSE_EXPORTER_LOGGING_URL` | Yes | | Cloud Foundry Log Stream URL |
| `logging.tls.ca`<br />`FIREHOSE_EXPORTER_LOGGING_TLS_CA` | No | | Path to ca cert to connect to rlp |
| `logging.tls.cert`<br />`FIREHOSE_EXPORTER_LOGGING_TLS_CERT` | Yes | | Path to cert to connect to rlp in mtls |
| `logging.tls.key`<br />`FIREHOSE_EXPORTER_LOGGING_TLS_KEY` | Yes | | Path to key to connect to rlp in mtls |
| `metrics.namespace`<br />`FIREHOSE_EXPORTER_METRICS_NAMESPACE` | No | `firehose` | Metrics Namespace |
| `metrics.environment`<br />`FIREHOSE_EXPORTER_METRICS_ENVIRONMENT` | Yes | | Environment label to be attached to metrics |
| `skip-ssl-verify`<br />`FIREHOSE_EXPORTER_SKIP_SSL_VERIFY` | No | `false` | Disable SSL Verify |
| `web.listen-address`<br />`FIREHOSE_EXPORTER_WEB_LISTEN_ADDRESS` | No | `:9186` | Address to listen on for web interface and telemetry |
| `web.telemetry-path`<br />`FIREHOSE_EXPORTER_WEB_TELEMETRY_PATH` | No | `/metrics` | Path under which to expose Prometheus metrics |
| `web.auth.username`<br />`FIREHOSE_EXPORTER_WEB_AUTH_USERNAME` | No | | Username for web interface basic auth |
| `web.auth.password`<br />`FIREHOSE_EXPORTER_WEB_AUTH_PASSWORD` | No | | Password for web interface basic auth |
| `web.tls.cert_file`<br />`FIREHOSE_EXPORTER_WEB_TLS_CERTFILE` | No | | Path to a file that contains the TLS certificate (PEM format). If the certificate is signed by a certificate authority, the file should be the concatenation of the server's certificate, any intermediates, and the CA's certificate |
| `web.tls.key_file`<br />`FIREHOSE_EXPORTER_WEB_TLS_KEYFILE` | No | | Path to a file that contains the TLS private key (PEM format) |
| `profiler.enable`<br />`FIREHOSE_EXPORTER_ENABLE_PROFILER` | No | `False` | Enable pprof profiling on app on /debug/pprof |
| `log.level`<br />`FIREHOSE_EXPORTER_LOG_LEVEL` | No | `info` | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] |
| `log.in_json`<br />`FIREHOSE_EXPORTER_LOG_IN_JSON` | No | `False` | Log in json |

### Metrics

For a list of [Cloud Foundry Firehose][firehose] metrics check the [Cloud Foundry Component Metrics][cfmetrics]
documentation.

The exporter returns additionally the following internal metrics:

| Metric | Description | Labels |
| ------ | ----------- | ------ |
| *
metrics.namespace*_total_envelopes_received | Total number of envelopes received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_envelope_received_timestamp | Number of seconds since 1970 since last envelope received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_total_metrics_received | Total number of metrics received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_metric_received_timestamp | Number of seconds since 1970 since last metric received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_total_container_metrics_received | Total number of container metrics received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_container_metric_received_timestamp | Number of seconds since 1970 since last container metric received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_total_counter_events_received | Total number of counter events received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_counter_event_received_timestamp | Number of seconds since 1970 since last counter event received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_total_http_received | Total number of http start stop received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_http_received_timestamp | Number of seconds since 1970 since last http start stop received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_total_value_metrics_received | Total number of value metrics received from Cloud Foundry Firehose | `environment` |
| *
metrics.namespace*_last_value_metric_received_timestamp | Number of seconds since 1970 since last value metric received from Cloud Foundry Firehose | `environment` |

## Contributing

Refer to [CONTRIBUTING.md](https://github.com/bosh-prometheus/firehose_exporter/blob/master/CONTRIBUTING.md).

## License

Apache License 2.0, see [LICENSE](https://github.com/bosh-prometheus/firehose_exporter/blob/master/LICENSE).

[cloudfoundry]: https://www.cloudfoundry.org/

[cfmetrics]: https://docs.cloudfoundry.org/loggregator/all_metrics.html

[faq]: https://github.com/bosh-prometheus/firehose_exporter/blob/master/FAQ.md

[firehose]: https://docs.cloudfoundry.org/loggregator/architecture.html#firehose

[golang]: https://golang.org/

[manifest]: https://github.com/bosh-prometheus/firehose_exporter/blob/master/manifest.yml

[prometheus]: https://prometheus.io/

[prometheus-boshrelease]: https://github.com/bosh-prometheus/prometheus-boshrelease
