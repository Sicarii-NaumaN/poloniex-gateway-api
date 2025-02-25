package decoder

import (
	"encoding/json"
	"io"
)

// DecodeJson helps to read json input and unpack it into go structure
func DecodeJson[T any](body io.ReadCloser) (res T, err error) {
	err = json.NewDecoder(body).Decode(&res)
	return res, err
}
