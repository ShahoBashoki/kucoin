package util

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Keyyer is an interface.
	Keyyer interface {
		// GetNode is a property.
		GetNode() string
	}

	key struct {
		node string
	}

	// Valuer is an interface.
	Valuer interface {
		// GetPayload is a property.
		GetPayload() map[string]any
		// GetB3 is a property.
		GetB3() string
	}

	value struct {
		payload map[string]any
		b3      string
	}

	// Recorder is an interface.
	Recorder interface {
		// GetKeyyer is a property.
		GetKeyyer() Keyyer
		// GetValuer is a property.
		GetValuer() Valuer
		// GetPartition is a property.
		GetPartition() int
	}

	// Redpandaer is an interface.
	Redpandaer interface {
		// RedpandaSend is a function.
		RedpandaSend(
			context.Context,
			map[string]any,
		) <-chan error
	}

	record struct {
		key       Keyyer
		value     Valuer
		partition int
	}

	// Recordser is an interface.
	Recordser interface {
		// GetRecorders is a property.
		GetRecorders() []Recorder
	}

	records struct {
		recorders []Recorder
	}

	redpanda struct {
		configConfigger  config.Configger
		logRuntimeLogger log.RuntimeLogger
		traceTracer      trace.Tracer
		utilUUIDer       UUIDer
	}
)

var (
	_ GetTracer            = (*redpanda)(nil)
	_ GetUUIDer            = (*redpanda)(nil)
	_ Keyyer               = (*key)(nil)
	_ Recorder             = (*record)(nil)
	_ Recordser            = (*records)(nil)
	_ Redpandaer           = (*redpanda)(nil)
	_ Valuer               = (*value)(nil)
	_ config.GetConfigger  = (*redpanda)(nil)
	_ json.Marshaler       = (*key)(nil)
	_ json.Marshaler       = (*record)(nil)
	_ json.Marshaler       = (*records)(nil)
	_ json.Marshaler       = (*value)(nil)
	_ log.GetRuntimeLogger = (*redpanda)(nil)
)

// NewKey is a function.
func NewKey(
	node string,
) *key {
	return &key{
		node: node,
	}
}

// NewValue is a function.
func NewValue(
	payload map[string]any,
	b3 string,
) *value {
	return &value{
		payload: payload,
		b3:      b3,
	}
}

// NewRecord is a function.
func NewRecord(
	key Keyyer,
	value Valuer,
	partition int,
) *record {
	return &record{
		key:       key,
		value:     value,
		partition: partition,
	}
}

// NewRecords is a function.
func NewRecords(
	recorders []Recorder,
) *records {
	return &records{
		recorders: recorders,
	}
}

// NewRedpanda is a function.
func NewRedpanda(
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer UUIDer,
) *redpanda {
	return &redpanda{
		configConfigger:  configConfigger,
		logRuntimeLogger: logRuntimeLogger,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
	}
}

// GetNode is a function.
func (key *key) GetNode() string {
	return key.node
}

// GetPayload is a function.
func (value *value) GetPayload() map[string]any {
	return value.payload
}

// GetB3 is a function.
func (value *value) GetB3() string {
	return value.b3
}

// GetKeyyer is a function.
func (record *record) GetKeyyer() Keyyer {
	return record.key
}

// GetValuer is a function.
func (record *record) GetValuer() Valuer {
	return record.value
}

// GetPartition is a function.
func (record *record) GetPartition() int {
	return record.partition
}

// GetRecorders is a function.
func (records *records) GetRecorders() []Recorder {
	return records.recorders
}

// GetMap is a function.
func (key *key) GetMap() map[string]any {
	return map[string]any{
		"node": key.GetNode(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (key *key) MarshalJSON() ([]byte, error) {
	return json.Marshal(key.GetMap())
}

// GetMap is a function.
func (value *value) GetMap() map[string]any {
	return map[string]any{
		"payload": value.GetPayload(),
		"b3":      value.GetB3(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (value *value) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.GetMap())
}

// GetMap is a function.
func (record *record) GetMap() map[string]any {
	return map[string]any{
		"key":       record.GetKeyyer(),
		"value":     record.GetValuer(),
		"partition": record.GetPartition(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (record *record) MarshalJSON() ([]byte, error) {
	return json.Marshal(record.GetMap())
}

// GetMap is a function.
func (records *records) GetMap() map[string]any {
	return map[string]any{
		"records": records.GetRecorders(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (records *records) MarshalJSON() ([]byte, error) {
	return json.Marshal(records.GetMap())
}

// GetConfigger is a function.
func (redpanda *redpanda) GetConfigger() config.Configger {
	return redpanda.configConfigger
}

// GetRuntimeLogger is a function.
func (redpanda *redpanda) GetRuntimeLogger() log.RuntimeLogger {
	return redpanda.logRuntimeLogger
}

// GetTracer is a function.
func (redpanda *redpanda) GetTracer() trace.Tracer {
	return redpanda.traceTracer
}

// GetUUIDer is a function.
func (redpanda *redpanda) GetUUIDer() UUIDer {
	return redpanda.utilUUIDer
}

// RedpandaSend is a function.
func (redpanda *redpanda) RedpandaSend(
	ctx context.Context,
	payload map[string]any,
) <-chan error {
	var traceSpan trace.Span

	ctx, traceSpan = redpanda.GetTracer().Start(
		ctx,
		"RedpandaSend",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	runtimeContext := NewRuntimeContext(ctx, redpanda.GetUUIDer())
	spanContext := NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":    "RedpandaSend",
		"rt_ctx":  runtimeContext,
		"sp_ctx":  spanContext,
		"config":  redpanda.GetConfigger(),
		"payload": payload,
	}

	redpanda.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	propagationMapCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, propagationMapCarrier)

	chanErr := make(chan error, 1)
	keyOne := NewKey(redpanda.GetConfigger().GetRuntimeConfigger().GetNode())
	valueOne := NewValue(payload, propagationMapCarrier.Get("b3"))
	partitionOne := 1
	recordOne := NewRecord(keyOne, valueOne, partitionOne)
	recordsOne := NewRecords([]Recorder{recordOne})

	body, err := recordsOne.MarshalJSON()
	if err != nil {
		redpanda.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrRecordsMarshalJSON.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrRecordsMarshalJSON.Error())

		go func() {
			chanErr <- err
		}()

		return chanErr
	}

	redpanda.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldBody, body).
		Debug(object.URIEmpty)

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		redpanda.GetConfigger().GetRedpandaConfigger().GetProxyURL(),
		bytes.NewReader(body),
	)
	if err != nil {
		redpanda.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrHTTPNewRequestWithContext.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrHTTPNewRequestWithContext.Error())

		go func() {
			chanErr <- err
		}()

		return chanErr
	}

	httpRequest.Header.Add(
		object.URIHTTPHeaderContentType,
		string(object.URIHTTPHeaderContentTypeAppKafka),
	)

	go func() {
		httpClient := &http.Client{
			Transport:     otelhttp.NewTransport(http.DefaultTransport),
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       object.NUMHTTPClientTimeout,
		}

		httpResponse, errHTTPClientDo := httpClient.Do(httpRequest)
		if errHTTPClientDo != nil {
			redpanda.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errHTTPClientDo).
				Error(object.ErrHTTPClientDo.Error())
			traceSpan.RecordError(errHTTPClientDo)
			traceSpan.SetStatus(codes.Error, object.ErrHTTPClientDo.Error())

			chanErr <- errHTTPClientDo

			return
		}

		redpanda.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldHTTPResponse, httpResponse).
			Debug(object.URIEmpty)

		if errHTTPClientDo = httpResponse.Body.Close(); errHTTPClientDo != nil {
			redpanda.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errHTTPClientDo).
				Error(object.ErrHTTPResponseBodyClose.Error())
			traceSpan.RecordError(errHTTPClientDo)
			traceSpan.SetStatus(codes.Error, object.ErrHTTPResponseBodyClose.Error())

			chanErr <- errHTTPClientDo

			return
		}

		chanErr <- nil
	}()

	return chanErr
}
