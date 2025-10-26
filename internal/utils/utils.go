package utils

import (
	"encoding/base64"
)

func EncodeURLSafeBase64(input []byte) string {
	return base64.RawURLEncoding.EncodeToString(input)
}
