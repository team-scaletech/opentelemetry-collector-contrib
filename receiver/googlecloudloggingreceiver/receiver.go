// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package googlecloudloggingreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudloggingreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

type googleLoggingReceiver struct {
	config       *Config
	settings     receiver.Settings
	logsConsumer consumer.Logs
	startTime    time.Time
	ctx          context.Context
	cancel       context.CancelFunc
	obsrecv      *receiverhelper.ObsReport
}

// newReceiver creates the Kubernetes events receiver with the given configuration.
func newReceiver(
	set receiver.Settings,
	config *Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	transport := "http"

	obsrecv, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             set.ID,
		Transport:              transport,
		ReceiverCreateSettings: set,
	})
	if err != nil {
		return nil, err
	}

	return &googleLoggingReceiver{
		settings:     set,
		config:       config,
		logsConsumer: consumer,
		startTime:    time.Now(),
		obsrecv:      obsrecv,
	}, nil
}

func (glr *googleLoggingReceiver) Start(ctx context.Context, _ component.Host) error {
	// Start receive logs
	return nil
}

func (glr *googleLoggingReceiver) Shutdown(context.Context) error {
	// Shutdown receive logs
	return nil
}
