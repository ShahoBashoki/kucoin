package object

import "errors"

var (
	// ErrBase64Decode2 is an error.
	ErrBase64Decode2 = errors.New("unrecognized level")
	// ErrGormOpen is an error.
	ErrGormOpen = errors.New("failed to open gorm")
	// ErrHTTPClientDo is an error.
	ErrHTTPClientDo = errors.New("failed to do http client")
	// ErrHTTPNewRequestWithContext is an error.
	ErrHTTPNewRequestWithContext = errors.New("failed to create a new http request with context")
	// ErrHTTPResponseBodyClose is an error.
	ErrHTTPResponseBodyClose = errors.New("failed to close http response body")
	// ErrJaegerNew is an error.
	ErrJaegerNew = errors.New("failed to create a jaeger exporter")
	// ErrRecordsMarshalJSON is an error.
	ErrRecordsMarshalJSON = errors.New("failed to marshall to byte array")
	// ErrRouterRun is an error.
	ErrRouterRun = errors.New("failed to router run")
	// ErrSDKResourceMerge is an error.
	ErrSDKResourceMerge = errors.New("failed to merge resources")
	// ErrSDKResourceNew is an error.
	ErrSDKResourceNew = errors.New("failed to create a new resource")
	// ErrSQL is an error.
	ErrSQL = errors.New("sql error")
	// ErrServerRun is an error.
	ErrServerRun = errors.New("failed to run http server")
	// ErrTracerProviderShutdown is an error.
	ErrTracerProviderShutdown = errors.New("failed to shutdown traceTracer provider")
)
