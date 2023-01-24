package socks

import (
	"bytes"
	"just/utils"
	"strings"
)

func ParseRawSocks(raw string) (map[string]string, error) {
	psIndex := strings.LastIndexByte(raw, '#')
	name := raw[psIndex+1:]
	raw = raw[:psIndex]
	dst, err := utils.Base64Decode(raw)
	if err != nil {
		return nil, err
	}

	authIndex := bytes.IndexByte(dst, ':')
	passIndex := bytes.IndexByte(dst, '@')
	hostIndex := bytes.LastIndexByte(dst, ':')

	var result = map[string]string{
		"name":     name,
		"cipher":   string(dst[:authIndex]),
		"password": string(dst[authIndex+1 : passIndex]),
		"server":   string(dst[passIndex+1 : hostIndex]),
		"port":     string(dst[hostIndex+1:]),
		"type":     "ss",
	}
	return result, nil
}
