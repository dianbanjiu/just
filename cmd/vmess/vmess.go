package vmess

import (
	"encoding/json"
	"fmt"
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

func ParseRawVmess(raw string) (map[string]string, error) {
	var vmess VmessInfo

	dst, err := utils.Base64Decode(raw)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dst, &vmess)
	if err != nil {
		return nil, err
	}
	var result = map[string]string{
		"name":    vmess.Name,
		"port":    vmess.Port,
		"server":  vmess.Server,
		"type":    "vmess",
		"uuid":    vmess.UUid,
		"cipher":  "auto",
		"alterId": fmt.Sprint(vmess.AlterId),
	}
	return result, nil
}
