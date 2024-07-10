// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlecloudmonitoringreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudmonitoringreceiver"

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const minCollectionIntervalSeconds = 60

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`

	Region            string    `mapstructure:"region"`
	ProjectID         string    `mapstructure:"project_id"`
	ServiceAccountKey string    `mapstructure:"service_account_key"`
	Services          []Service `mapstructure:"services"`
}

type Service struct {
	ServiceName string        `mapstructure:"service_name"`
	Delay       time.Duration `mapstructure:"delay"`
	MetricName  string        `mapstructure:"metric_name"`
}

func (config *Config) Validate() error {
	if config.CollectionInterval.Seconds() < minCollectionIntervalSeconds {
		return fmt.Errorf("\"collection_interval\" must be not lower than %v seconds, current value is %v seconds", minCollectionIntervalSeconds, config.CollectionInterval.Seconds())
	}

	if len(config.Services) == 0 {
		return errors.New("missing required field \"services\" or its value is empty")
	}

	for _, service := range config.Services {
		if err := service.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (service Service) Validate() error {
	if service.ServiceName == "" {
		return errors.New("field \"service_name\" is required and cannot be empty for service configuration")
	}

	if service.Delay < 0 {
		return errors.New("field \"delay\" cannot be negative for service configuration")
	}

	return nil
}
