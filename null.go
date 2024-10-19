package query

import (
	"bytes"
	"encoding/json"
)

var jsonNull = []byte("null")

type Null[T any] struct {
	Value T
	Valid bool
}

func NewNull[T any](value T, valid bool) Null[T] {
	return Null[T]{
		Value: value,
		Valid: valid,
	}
}

func (n Null[T]) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Value)
	}

	return jsonNull, nil
}

func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		*n = Null[T]{}
		return nil
	}

	err := json.Unmarshal(data, &n.Value)
	n.Valid = err == nil

	return err
}
