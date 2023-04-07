package dao

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"

	"github.com/ShahoBashoki/kucoin/object"
	"gorm.io/gorm"
)

type (
	// Cursorer is an interface.
	Cursorer interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
		// GetOffset is a function.
		GetOffset() uint32
		// Query is a function.
		Query(
			table string,
		) func(*gorm.DB) *gorm.DB
	}

	// GetCursorer is an interface.
	GetCursorer interface {
		// GetCursorer is a function.
		GetCursorer() Cursorer
	}

	cursor struct {
		offset uint32
	}
)

var (
	_ Cursorer                   = (*cursor)(nil)
	_ encoding.BinaryMarshaler   = (*cursor)(nil)
	_ encoding.BinaryUnmarshaler = (*cursor)(nil)
	_ json.Marshaler             = (*cursor)(nil)
	_ object.GetMap              = (*cursor)(nil)
)

// NewCursor is a function.
func NewCursor(
	offset uint32,
) *cursor {
	return &cursor{
		offset: offset,
	}
}

// CursorerComparer is a function.
func CursorerComparer(
	first Cursorer,
	second Cursorer,
) bool {
	return first.GetOffset() == second.GetOffset()
}

// GetOffset is a function.
func (cursor *cursor) GetOffset() uint32 {
	return cursor.offset
}

// Query is a function.
func (cursor *cursor) Query(
	_ string,
) func(*gorm.DB) *gorm.DB {
	return func(
		gormDB *gorm.DB,
	) *gorm.DB {
		gormDB.Offset(int(cursor.GetOffset()))

		return gormDB
	}
}

// GetMap is a function.
func (cursor *cursor) GetMap() map[string]any {
	return map[string]any{
		"offset": cursor.GetOffset(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (cursor *cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(cursor.GetMap())
}

// MarshalBinary is a function.
// read more https://pkg.go.dev/encoding#BinaryMarshaler
func (cursor *cursor) MarshalBinary() ([]byte, error) {
	var bytesBuffer bytes.Buffer
	_, err := fmt.Fprintln(&bytesBuffer, cursor.offset)

	return bytesBuffer.Bytes(), err
}

// UnmarshalBinary is a function.
// read more https://pkg.go.dev/encoding#BinaryUnmarshaler
func (cursor *cursor) UnmarshalBinary(
	data []byte,
) error {
	bytesBuffer := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(bytesBuffer, &cursor.offset)

	return err
}
