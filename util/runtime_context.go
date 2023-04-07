package util

import (
	"context"
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
	"github.com/google/uuid"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/metadata"
)

type (
	// RuntimeContexter is an interface.
	RuntimeContexter interface {
		// GetMetadata is a function.
		GetMetadata() metadata.MD
		// GetClientHost is a function.
		GetClientHost() string
		// GetClientPort is a function.
		GetClientPort() string
		// GetUserID is a function.
		GetUserID() uuid.UUID
	}

	runtimeContext struct {
		md         metadata.MD
		clientHost string
		clientPort string
		userID     uuid.UUID
	}
)

var (
	_ RuntimeContexter = (*runtimeContext)(nil)
	_ json.Marshaler   = (*runtimeContext)(nil)
)

// NewRuntimeContext is a function.
func NewRuntimeContext(
	ctx context.Context,
	utilUUIDer UUIDer,
) *runtimeContext {
	values := grpcTags.Extract(ctx).Values()

	metadataMD, ok := values[object.URIRuntimeContextMetadata].(metadata.MD)
	if !ok {
		metadataMD = metadata.MD{}
	}

	clientHost, ok := values[object.URIRuntimeContextClientHost].(string)
	if !ok {
		clientHost = object.URIEmpty
	}

	clientPort, ok := values[object.URIRuntimeContextClientPort].(string)
	if !ok {
		clientPort = object.URIEmpty
	}

	userID, ok := values[object.URIRuntimeContextUserID].(string)
	if !ok {
		userID = uuid.Nil.String()
	}

	userUUID, err := utilUUIDer.Parse(userID)
	if err != nil {
		userUUID = uuid.Nil
	}

	runtimeContext := &runtimeContext{
		md:         metadataMD,
		clientHost: clientHost,
		clientPort: clientPort,
		userID:     userUUID,
	}

	return runtimeContext
}

// GetMetadata is a function.
func (runtimeContext *runtimeContext) GetMetadata() metadata.MD {
	return runtimeContext.md
}

// GetClientHost is a function.
func (runtimeContext *runtimeContext) GetClientHost() string {
	return runtimeContext.clientHost
}

// GetClientPort is a function.
func (runtimeContext *runtimeContext) GetClientPort() string {
	return runtimeContext.clientPort
}

// GetUserID is a function.
func (runtimeContext *runtimeContext) GetUserID() uuid.UUID {
	return runtimeContext.userID
}

// GetMap is a function.
func (runtimeContext *runtimeContext) GetMap() map[string]any {
	return map[string]any{
		"md":          runtimeContext.GetMetadata(),
		"client_host": runtimeContext.GetClientHost(),
		"client_port": runtimeContext.GetClientPort(),
		"user_id":     runtimeContext.GetUserID(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (runtimeContext *runtimeContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(runtimeContext.GetMap())
}
