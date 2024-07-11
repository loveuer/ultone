package cache

import (
	"encoding/json"
	"ultone/internal/interfaces"
)

var (
	Client interfaces.Cacher
)

type encoded_value interface {
	MarshalBinary() ([]byte, error)
}

type decoded_value interface {
	UnmarshalBinary(bs []byte) error
}

func handleValue(value any) ([]byte, error) {
	var (
		bs  []byte
		err error
	)

	if imp, ok := value.(encoded_value); ok {
		bs, err = imp.MarshalBinary()
	} else {
		bs, err = json.Marshal(value)
	}

	return bs, err
}
