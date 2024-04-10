package openrtb

import (
	"encoding/json"
	"errors"
	"strconv"
)

// NumberOrString attempts to fix OpenRTB incompatibilities
// of exchanges. On decoding, it can handle numbers and strings.
// On encoding, it will generate a number, as intended by the
// standard.
type NumberOrString int

// UnmarshalJSON implements json.Unmarshaler
func (n *NumberOrString) UnmarshalJSON(data []byte) (err error) {
	var v int

	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		err = json.Unmarshal(data[1:len(data)-1], &v)
	} else {
		err = json.Unmarshal(data, &v)
	}
	if err != nil {
		return err
	}

	*n = NumberOrString(v)
	return nil
}

// StringOrNumber attempts to fix OpenRTB incompatibilities
// of exchanges. On decoding, it can handle numbers and strings.
// On encoding, it will generate a string, as intended by the
// standard.
type StringOrNumber string

// UnmarshalJSON implements json.Unmarshaler
func (n *StringOrNumber) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		var v string
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*n = StringOrNumber(v)
	} else {
		var v int
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*n = StringOrNumber(strconv.Itoa(v))
	}
	return nil
}

// BoolOrNumber attemps to fix OpenRTB incompatibilities where a field is expected as bool but the spec expects int values. This was not seen till now, but some mediation partners will follow the spec closely.
type BoolOrNumber bool

// UnmarshalJSON implements json.Unmarshaler
func (b *BoolOrNumber) UnmarshalJSON(data []byte) error {
	var val interface{}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch v := val.(type) {
	case bool:
		*b = BoolOrNumber(v)
	case float64:
		// When unmarshaling JSON into an interface value, Unmarshal stores JSON numbers in the interface value float64
		*b = BoolOrNumber(v != 0)
	default:

		return errors.New("BoolOrInt: invalid type")
	}

	return nil
}
