package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/sagikazarmark/kcconf/api"
)

func main() {
	// TODO: make logger configurable
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	kafkaConnectURL := pflag.String("kafka-connect-url", envDefault("KAFKA_CONNECT_URL", "http://127.0.0.1:8083"), "Kafka Connect REST API URL")
	connectorsFile := pflag.String("connectors-file", envDefault("CONNECTORS_FILE", "/etc/kcconf/connectors.yaml"), "Config file with a list of connectors")

	pflag.Parse()

	if *kafkaConnectURL == "" {
		logger.Error("Kafka Connect URL is required")
	}

	if *connectorsFile == "" {
		logger.Error("Connectors file is required")
	}

	client, err := api.NewClient(*kafkaConnectURL)
	if err != nil {
		panic(err)
	}

	logger.Info("connecting to Kafka Connect", slog.String("url", *kafkaConnectURL))

	err = backoff.RetryNotify(
		func() error {
			resp, err := client.ServerInfo(context.Background())
			if err != nil {
				return err
			}

			resp.Body.Close()

			return nil
		},
		backoff.NewExponentialBackOff(),
		func(err error, d time.Duration) {
			logger.Debug("Kafka Connect is unavailable, retrying", slog.String("url", *kafkaConnectURL), slog.Any("error", err), slog.Duration("duration", d))
		},
	)
	if err != nil {
		panic(err)
	}

	logger.Info("loading connector configuration")

	file, err := os.Open("connectors.yaml")
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(file)

	type connectorConfig struct {
		Name   string            `yaml:"name"`
		Config map[string]string `yaml:"config"`
	}

	var configs []connectorConfig

	err = decoder.Decode(&configs)
	if err != nil {
		panic(err)
	}

	for _, config := range configs {
		logger.Info("configuring connector", slog.String("name", config.Name))

		resp, err := client.PutConnectorConfig(context.Background(), config.Name, config.Config)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == http.StatusCreated {
			logger.Info("connector created", slog.String("name", config.Name))
		} else if resp.StatusCode == http.StatusOK {
			logger.Info("connector updated", slog.String("name", config.Name))
		} else if resp.StatusCode == http.StatusConflict {
			// TODO: retry
			logger.Warn("connector creation could not complete: rebalance in progress", slog.String("name", config.Name))
		} else if resp.StatusCode > 299 {
			type errorResponse struct {
				ErrorCode int    `json:"error_code"`
				Message   string `json:"message"`
			}

			var apiResp errorResponse

			decoder := json.NewDecoder(resp.Body)

			err := decoder.Decode(&apiResp)
			if err != nil {
				logger.Error("connector create/update failed, but the error response is gibberish", slog.Int("status-code", resp.StatusCode))
			}

			logger.Error(fmt.Sprintf("connector create/update failed: %s", apiResp.Message), slog.Int("error-code", apiResp.ErrorCode))
		}

		resp.Body.Close()
	}
}

func envDefault(envVar string, def string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}

	return def
}
