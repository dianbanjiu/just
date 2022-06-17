package vmess

import (
	"encoding/json"
	"just/utils"
)

type VmessInfo struct {
	Name    string `json:"ps" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	Server  string `json:"add" yaml:"server"`
	Port    string `json:"port" yaml:"port"`
	UUid    string `json:"id" yaml:"uuid"`
	AlterId int    `json:"aid" yaml:"alterId"`
	Cipher  string `json:"cipher" yaml:"cipher"`
}

func ParseRawVmess(raw string) (VmessInfo, error) {
	var vmess VmessInfo

	dst, err := utils.Base64Decode(raw)

	err = json.Unmarshal(dst, &vmess)
	if err != nil {
		return vmess, err
	}
	vmess.Type = "vmess"
	vmess.Cipher = "auto"
	return vmess, nil
}
