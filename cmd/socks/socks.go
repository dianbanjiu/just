package socks

import (
	"bytes"
	"just/utils"
	"strings"
)

type SocksInfo struct {
	Name     string `json:"ps" yaml:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Server   string `json:"server" yaml:"server"`
	Port     string `json:"port" yaml:"port"`
	Cipher   string `json:"cipher" yaml:"cipher"`
	Password string `json:"password" yaml:"password"`
}

func ParseRawSocks(raw string) (SocksInfo, error) {
	var socks SocksInfo
	psIndex := strings.LastIndexByte(raw, '#')
	socks.Name = raw[psIndex+1:]

	raw = raw[:psIndex]
	dst, err := utils.Base64Decode(raw)
	if err != nil {
		return socks, err
	}

	authIndex := bytes.IndexByte(dst, ':')
	socks.Cipher = string(dst[:authIndex])
	passIndex := bytes.IndexByte(dst, '@')
	socks.Password = string(dst[authIndex+1 : passIndex])
	hostIndex := bytes.LastIndexByte(dst, ':')
	socks.Server = string(dst[passIndex+1 : hostIndex])
	socks.Port = string(dst[hostIndex+1:])
	socks.Type = "ss"
	return socks, nil
}
