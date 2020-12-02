package logstream

import (
	"net/http"

	"code.cloudfoundry.org/go-loggregator"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/firehose_exporter/metrics"
)

type LogStream struct {
	url               string
	caPath            string
	certPath          string
	keyPath           string
	useRLPGateway     bool
	skipSSLValidation bool
	subscriptionID    string
	metricsStore      *metrics.Store
	messages          <-chan *events.Envelope
	consumer          *V2Adapter
	httpClient        doer
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(
	url string,
	caPath string,
	certPath string,
	keyPath string,
	useRLPGateway bool,
	skipSSLValidation bool,
	subscriptionID string,
	metricsStore *metrics.Store,
	httpClient doer,
) *LogStream {
	return &LogStream{
		url:               url,
		caPath:            caPath,
		certPath:          certPath,
		keyPath:           keyPath,
		useRLPGateway:     useRLPGateway,
		skipSSLValidation: skipSSLValidation,
		subscriptionID:    subscriptionID,
		metricsStore:      metricsStore,
		messages:          make(<-chan *events.Envelope),
		httpClient:        httpClient,
	}
}

// Start processes both errors and messages until both channels are closed
// It then closes the underlying consumer.
func (n *LogStream) Start() error {
	log.Info("Starting Firehose Nozzle...")
	defer log.Info("Firehose Nozzle shutting down...")
	err := n.consumeLogstream()
	if err != nil {
		return err
	}
	n.parseEnvelopes()
	return nil
}

func (n *LogStream) consumeLogstream() error {
	streamer, err := n.buildStreamer()
	if err != nil {
		return err
	}
	a := NewV2Adapter(streamer)
	n.messages = a.Firehose(n.subscriptionID)
	return nil
}

func (n *LogStream) buildStreamer() (Streamer, error) {
	if n.useRLPGateway {
		return loggregator.NewRLPGatewayClient(
			n.url,
			loggregator.WithRLPGatewayHTTPClient(n.httpClient),
		), nil
	}
	loggregatorTLSConfig, err := loggregator.NewEgressTLSConfig(n.caPath, n.certPath, n.keyPath)
	if err != nil {
		return nil, err
	}
	loggregatorTLSConfig.InsecureSkipVerify = n.skipSSLValidation
	return loggregator.NewEnvelopeStreamConnector(
		n.url,
		loggregatorTLSConfig,
		loggregator.WithEnvelopeStreamBuffer(10000, func(missed int) {
			log.Infof("dropped %d envelope batches", missed)
		}),
	), nil
}

// parseEnvelopes will read and process both errs and messages, until
// both are closed, at which time it will close the consumer and return
func (n *LogStream) parseEnvelopes() {
	defer n.consumer.Close()

	for {
		select {
		case envelope, ok := <-n.messages:
			if !ok {
				return
			}
			n.metricsStore.AddMetric(envelope)
		}
	}
}
