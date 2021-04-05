package toolkit

import (
	"encoding/json"
	"strings"
)

func UnmarshalStrict(s string, v interface{}) error {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
