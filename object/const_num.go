package object

import "time"

const (
	// NUMHTTPClientTimeout is a variable.
	NUMHTTPClientTimeout = 1500 * time.Millisecond
	// NUMLogConfigDefaultLogMaxSize is a variable.
	NUMLogConfigDefaultLogMaxSize = 100
	// NUMRuntimeConfigDefaultRuntimeKucoinPaginationRequestSize is a variable.
	NUMRuntimeConfigDefaultRuntimeKucoinPaginationRequestSize = 500
	// NUMSystemGracefulShutdown is a variable.
	NUMSystemGracefulShutdown = 5 * time.Second
	// NUMTopTickerChangeRateCount is a variable.
	NUMTopTickerChangeRateCount = 10
)
