package util

import "github.com/google/uuid"

type (
	// UUIDer is an interface.
	UUIDer interface {
		// NewRandom is a function.
		NewRandom() (uuid.UUID, error)
		// Parse is a function.
		Parse(
			string,
		) (uuid.UUID, error)
	}

	// GetUUIDer is an interface.
	GetUUIDer interface {
		// GetUUIDer is a function.
		GetUUIDer() UUIDer
	}

	uuidV2 struct {
		uuid.UUID
	}
)

var _ UUIDer = (*uuidV2)(nil)

// NewUUID is a function.
func NewUUID() *uuidV2 {
	return &uuidV2{
		UUID: uuid.Nil,
	}
}

// NewRandom is a function.
func (*uuidV2) NewRandom() (uuid.UUID, error) {
	return uuid.NewRandom()
}

// Parse is a function.
func (*uuidV2) Parse(
	s string,
) (uuid.UUID, error) {
	return uuid.Parse(s)
}
