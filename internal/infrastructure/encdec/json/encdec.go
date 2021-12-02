package json

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type EncDec struct{}

func (e *EncDec) Encode(in interface{}) ([]byte, error) {
	d, err := json.Marshal(in)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return d, nil
}
func (e *EncDec) Decode(data []byte, in interface{}) error {
	return errors.WithStack(json.Unmarshal(data, in))
}
