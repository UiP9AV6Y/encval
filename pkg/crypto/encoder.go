package crypto

import (
	"encoding/base64"
)

var encoder = base64.StdEncoding.Strict()

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
func Decode(data Data) ([]byte, error) {
	result := make([]byte, encoder.DecodedLen(len(data)))
	n, err := encoder.Decode(result, []byte(data))
	if err != nil {
		return nil, err
	}

	return result[:n], nil
}
