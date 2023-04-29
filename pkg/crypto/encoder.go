package crypto

import (
	"encoding/base64"
)

var encoder = base64.RawStdEncoding.Strict()

// Encode is an opiniated Base64 encoder function
// using a strict variant of base64.RawStdEncoding
// for data processing
func Encode(data []byte) (Data, error) {
	result := make([]byte, encoder.EncodedLen(len(data)))
	encoder.Encode(result, data)

	return Data(result), nil
}

// Decode is an opiniated Base64 decoder function
// using a strict variant of base64.RawStdEncoding
// for data processing
func Decode(data []byte) (Data, error) {
	result := make([]byte, encoder.DecodedLen(len(data)))
	n, err := encoder.Decode(result, data)
	if err != nil {
		return nil, err
	}

	return Data(result[:n]), nil
}
