package utils

import (
	"encoding/base64"
	"strings"
)

func Base64Decode(raw string) ([]byte, error) {
	if len(raw)%4 != 0 {
		raw += strings.Repeat("=", 4-len(raw)%4)
	}
	return base64.StdEncoding.DecodeString(raw)
}
