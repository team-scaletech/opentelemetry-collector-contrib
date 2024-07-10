package googlecloudloggingreceiver

import (
	"errors"
	"fmt"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	minCollectionIntervalSeconds = 60
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`

	Region            string     `mapstructure:"region"`
	ProjectID         string     `mapstructure:"project_id"`
	ServiceAccountKey string     `mapstructure:"service_account_key"`
	Resources         []Resource `mapstructure:"resources"`
}

type Resource struct {
	ResourceType      string             `mapstructure:"resource_type"`
	AdditionalFilters []AdditionalFilter `mapstructure:"additionalFilters"`
	ResourceNames     []string           `mapstructure:"resource_names"`
}

type AdditionalFilter struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (config *Config) Validate() error {
	if config.CollectionInterval.Seconds() < minCollectionIntervalSeconds {
		return fmt.Errorf("\"collection_interval\" must be not lower than %v seconds, current value is %v seconds", minCollectionIntervalSeconds, config.CollectionInterval.Seconds())
	}

	if config.Region == "" {
		return errors.New("field \"region\" is required and cannot be empty")
	}

	if config.ProjectID == "" {
		return errors.New("field \"project_id\" is required and cannot be empty")
	}

	if config.ServiceAccountKey == "" {
		return errors.New("field \"service_account_key\" is required and cannot be empty")
	}

	if len(config.Resources) == 0 {
		return errors.New("field \"resources\" is required and cannot be empty")
	}

	for _, resource := range config.Resources {
		if err := resource.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (resource Resource) Validate() error {
	if resource.ResourceType == "" {
		return errors.New("field \"resource_type\" is required and cannot be empty")
	}

	for _, filter := range resource.AdditionalFilters {
		if err := filter.Validate(); err != nil {
			return err
		}
	}

	if len(resource.ResourceNames) == 0 {
		return errors.New("field \"resource_names\" is required and cannot be empty")
	}

	for _, name := range resource.ResourceNames {
		if name == "" {
			return errors.New("field \"resource_names\" contains empty resource name")
		}
	}

	return nil
}

func (filter AdditionalFilter) Validate() error {
	if filter.Name == "" {
		return errors.New("field \"name\" is required and cannot be empty")
	}

	if filter.Value == "" {
		return errors.New("field \"value\" is required and cannot be empty")
	}

	return nil
}
