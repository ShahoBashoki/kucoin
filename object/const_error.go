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
	// ErrKucoinServiceReadPaginationData is an error.
	ErrKucoinServiceReadPaginationData = errors.New("failed to kucoin service read pagination data")
	// ErrOrderBookKucoinServiceGetList is an error.
	ErrOrderBookKucoinServiceGetList = errors.New("failed to order book kucoin service get list")
	// ErrOrderBookServiceGetListFromRepository is an error.
	ErrOrderBookServiceGetListFromRepository = errors.New(
		"failed to order book service get list from repository",
	)
	// ErrOrderKucoinServiceGetList is an error.
	ErrOrderKucoinServiceGetList = errors.New("failed to order kucoin service get list")
	// ErrOrderRepositoryCreate is an error.
	ErrOrderRepositoryCreate = errors.New("failed to order repository create")
	// ErrOrderRepositoryDelete is an error.
	ErrOrderRepositoryDelete = errors.New("failed to order repository delete")
	// ErrOrderRepositoryDeleteAll is an error.
	ErrOrderRepositoryDeleteAll = errors.New("failed to order repository delete all")
	// ErrOrderRepositoryRead is an error.
	ErrOrderRepositoryRead = errors.New("failed to order repository read")
	// ErrOrderRepositoryReadList is an error.
	ErrOrderRepositoryReadList = errors.New("failed to order repository read list")
	// ErrOrderRepositoryUpdate is an error.
	ErrOrderRepositoryUpdate = errors.New("failed to order repository update")
	// ErrOrderServiceCreate is an error.
	ErrOrderServiceCreate = errors.New("failed to order service create")
	// ErrOrderServiceDeleteAll is an error.
	ErrOrderServiceDeleteAll = errors.New("failed to order service delete all")
	// ErrOrderServiceGetListFromRemote is an error.
	ErrOrderServiceGetListFromRemote = errors.New("failed to order service get list from remote")
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
	// ErrSTRCONVParseFloat is an error.
	ErrSTRCONVParseFloat = errors.New("failed to strconv parse float")
	// ErrServerRun is an error.
	ErrServerRun = errors.New("failed to run http server")
	// ErrTickerKucoinServiceGetList is an error.
	ErrTickerKucoinServiceGetList = errors.New("failed to ticker kucoin service get list")
	// ErrTickerRepositoryCreate is an error.
	ErrTickerRepositoryCreate = errors.New("failed to ticker repository create")
	// ErrTickerRepositoryDelete is an error.
	ErrTickerRepositoryDelete = errors.New("failed to ticker repository delete")
	// ErrTickerRepositoryDeleteAll is an error.
	ErrTickerRepositoryDeleteAll = errors.New("failed to ticker repository delete all")
	// ErrTickerRepositoryRead is an error.
	ErrTickerRepositoryRead = errors.New("failed to ticker repository read")
	// ErrTickerRepositoryReadList is an error.
	ErrTickerRepositoryReadList = errors.New("failed to ticker repository read list")
	// ErrTickerRepositoryUpdate is an error.
	ErrTickerRepositoryUpdate = errors.New("failed to ticker repository update")
	// ErrTickerServiceCreate is an error.
	ErrTickerServiceCreate = errors.New("failed to ticker service create")
	// ErrTickerServiceDeleteAll is an error.
	ErrTickerServiceDeleteAll = errors.New("failed to ticker service delete all")
	// ErrTickerServiceGetListFromRepository is an error.
	ErrTickerServiceGetListFromRepository = errors.New(
		"failed to ticker service get list from repository",
	)
	// ErrTickerServiceGetListFromRemote is an error.
	ErrTickerServiceGetListFromRemote = errors.New("failed to ticker service get list from remote")
	// ErrTracerProviderShutdown is an error.
	ErrTracerProviderShutdown = errors.New("failed to shutdown traceTracer provider")
	// ErrTypeAssertion is an error.
	ErrTypeAssertion = errors.New("failed to assert type")
	// ErrUUIDerNewRandom is an error.
	ErrUUIDerNewRandom = errors.New("failed to new random uuid")
	// ErrUUIDerParse is an error.
	ErrUUIDerParse = errors.New("failed to parse uuid")
	// ErrUnmarshalJSON is an error.
	ErrUnmarshalJSON = errors.New("failed to unmarshal json")
)
