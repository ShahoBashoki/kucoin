package om

import "github.com/google/uuid"

// OMer is an interface.
type OMer interface {
	// GetID is a function.
	GetID() uuid.UUID
}

// OMerComparer is a function.
func OMerComparer(
	first OMer,
	second OMer,
) bool {
	return first.GetID() == second.GetID()
}
