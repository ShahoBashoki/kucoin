package log

import (
	"bytes"
	"io"
	stdLog "log"
	"strings"

	"github.com/ShahoBashoki/kucoin/object"
	"go.uber.org/zap"
)

type redirectStdLogWriter struct {
	zapLog *zap.Logger
}

var _ io.Writer = (*redirectStdLogWriter)(nil)

// Write is a function.
func (redirectStdLogWriter *redirectStdLogWriter) Write(
	p []byte,
) (int, error) {
	str := string(bytes.TrimSpace(p))
	if strings.HasPrefix(str, "http: panic serving") {
		redirectStdLogWriter.zapLog.Error(str)
	} else {
		redirectStdLogWriter.zapLog.Info(str)
	}

	return len(str), nil
}

// RedirectStdLog is a function.
func RedirectStdLog(
	zapLogger *zap.Logger,
) {
	stdLog.SetFlags(0)
	stdLog.SetPrefix(object.URIEmpty)

	newZapLogger := zapLogger.WithOptions(zap.AddCallerSkip(3))

	stdLog.SetOutput(&redirectStdLogWriter{
		newZapLogger,
	})
}
