# kcconf

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/kcconf/ci.yaml?style=flat-square)](https://github.com/sagikazarmark/kcconf/actions/workflows/ci.yaml)
[![built with nix](https://img.shields.io/badge/builtwith-nix-7d81f7?style=flat-square)](https://builtwithnix.org)

**[Kafka Connect](https://docs.confluent.io/platform/current/connect/index.html) configurator tool.**

## Usage

This tool allows configuring Kafka Connect when running in [distributed mode](https://docs.confluent.io/platform/current/connect/index.html#distributed-workers)
(which is the default mode of operation when running basically any of the containerized versions).

It allows configuring a list of connectors and keeping those configurations up-to-date.

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

TODO

## License

The project is licensed under the [MIT License](LICENSE).
