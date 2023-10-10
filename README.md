# kcconf

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/kcconf/ci.yaml?style=flat-square)](https://github.com/sagikazarmark/kcconf/actions/workflows/ci.yaml)
[![built with nix](https://img.shields.io/badge/builtwith-nix-7d81f7?style=flat-square)](https://builtwithnix.org)

**[Kafka Connect](https://docs.confluent.io/platform/current/connect/index.html) configurator tool.**

## Usage

This tool allows configuring Kafka Connect when running in [distributed mode](https://docs.confluent.io/platform/current/connect/index.html#distributed-workers)
(which is the default mode of operation when running basically any of the containerized versions).

It allows configuring a list of connectors and keeping those configurations up-to-date.

Create a YAML file with a list of connectors:

```yaml
- name: mock
  config:
    connector.class: org.apache.kafka.connect.tools.MockSinkConnector
    key.converter: "org.apache.kafka.connect.storage.StringConverter"
    value.converter: "org.apache.kafka.connect.json.JsonConverter"
    value.converter.schemas.enable: "false"
    schemas.enable: "false"
    topics.regex: "^om_[A-Za-z0-9]+(?:_[A-Za-z0-9]+)*_events$"
    errors.tolerance: "all"
    errors.retry.timeout: "30"
```

Run `kcconf`:

```shell
kcconf --kafka-connect-url http://127.0.0.1:8080 --connectors-file connectors.yaml
```

**Pro tip:** You can run `kcconf` in a Docker Compose setup:

```yaml
version: "3.9"

services:
  # ...

  kcconf:
    image: ghcr.io/sagikazarmark/kcconf
    depends_on:
      # This is technically not necessary, because kcconf will retry connecting to Kafka Connect,
      # but if you have a health check it doesn't hurt either.
      kafka-connect:
        condition: service_healthy
    environment:
      KAFKA_CONNECT_URL: http://kafka-connect:8080
    volumes:
      - $PWD/connectors.yaml:/etc/kcconf/connectors.yaml
```

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

TODO

## License

The project is licensed under the [MIT License](LICENSE).
