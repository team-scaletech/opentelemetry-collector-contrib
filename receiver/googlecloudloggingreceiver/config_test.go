// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlecloudloggingreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudloggingreceiver"

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudloggingreceiver/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(metadata.Type, "").String())
	require.NoError(t, err)
	require.NoError(t, sub.Unmarshal(cfg))

	assert.Equal(t,
		&Config{
			ControllerConfig: scraperhelper.ControllerConfig{
				CollectionInterval: 120 * time.Second,
			},
			Region:            "us-central1",
			ProjectID:         "my-project-id",
			ServiceAccountKey: "path/to/service_account.json",
			Resources: []Resource{
				{
					ResourceType: "cloudfunctions",
					AdditionalFilters: []AdditionalFilter{
						{
							Name:  "resource.labels.service_name",
							Value: "service-name",
						},
						{
							Name:  "resource.labels.location",
							Value: "us-east-1",
						},
					},
					ResourceNames: []string{
						"projects/my-project-id1",
						"projects/my-project-id2",
					},
				},
				{
					ResourceType: "another-resource-type",
					ResourceNames: []string{
						"projects/my-project-id3",
						"projects/my-project-id4",
					},
				},
			},
		},
		cfg,
	)
}

func TestValidateAdditionalFilter(t *testing.T) {
	testCases := map[string]struct {
		name         string
		value        string
		requireError bool
	}{
		"All required fields are populated": {"resource.labels.service_name", "service-name", false},
		"No name":                           {"", "service-name", true},
		"No value":                          {"resource.labels.service_name", "", true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			filter := AdditionalFilter{
				Name:  testCase.name,
				Value: testCase.value,
			}

			err := filter.Validate()

			if testCase.requireError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateResource(t *testing.T) {
	testCases := map[string]struct {
		resourceType      string
		additionalFilters []AdditionalFilter
		resourceNames     []string
		requireError      bool
	}{
		"All required fields are populated":   {"cloudfunctions", []AdditionalFilter{{"resource.labels.service_name", "service-name"}}, []string{"projects/my-project-id1"}, false},
		"No resource type":                    {"", []AdditionalFilter{{"resource.labels.service_name", "service-name"}}, []string{"projects/my-project-id1"}, true},
		"Empty resource names":                {"cloudfunctions", []AdditionalFilter{{"resource.labels.service_name", "service-name"}}, nil, true},
		"Resource names contain empty string": {"cloudfunctions", []AdditionalFilter{{"resource.labels.service_name", "service-name"}}, []string{""}, true},
		"Invalid additional filter":           {"cloudfunctions", []AdditionalFilter{{"", "service-name"}}, []string{"projects/my-project-id1"}, true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource := Resource{
				ResourceType:      testCase.resourceType,
				AdditionalFilters: testCase.additionalFilters,
				ResourceNames:     testCase.resourceNames,
			}

			err := resource.Validate()

			if testCase.requireError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	resource := Resource{
		ResourceType: "cloudfunctions",
		AdditionalFilters: []AdditionalFilter{
			{
				Name:  "resource.labels.service_name",
				Value: "service-name",
			},
		},
		ResourceNames: []string{"projects/my-project-id1"},
	}

	testCases := map[string]struct {
		region            string
		projectID         string
		serviceAccountKey string
		resources         []Resource
		requireError      bool
	}{
		"All required fields are populated": {"us-central1", "my-project-id", "path/to/service_account.json", []Resource{resource}, false},
		"No region":                         {"", "my-project-id", "path/to/service_account.json", []Resource{resource}, true},
		"No project ID":                     {"us-central1", "", "path/to/service_account.json", []Resource{resource}, true},
		"No service account key":            {"us-central1", "my-project-id", "", []Resource{resource}, true},
		"No resources":                      {"us-central1", "my-project-id", "path/to/service_account.json", nil, true},
		"Invalid resource in resources":     {"us-central1", "my-project-id", "path/to/service_account.json", []Resource{{}}, true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			cfg := &Config{
				ControllerConfig: scraperhelper.ControllerConfig{
					CollectionInterval: 120 * time.Second,
				},
				Region:            testCase.region,
				ProjectID:         testCase.projectID,
				ServiceAccountKey: testCase.serviceAccountKey,
				Resources:         testCase.resources,
			}

			err := cfg.Validate()

			if testCase.requireError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
