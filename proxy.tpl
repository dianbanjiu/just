{
    "inbounds": [
        {
            "port": 1006,
            "listen": "127.0.0.1",
            "protocol": "socks",
            "settings": {
                "udp": true
            }
        }
    ],
    "outbounds": [
        {
            "protocol": "vmess",
            "settings": {
                "vnext": [{{$i:=0}}
                    {{range .}}{{if eq .type "vmess"}}{{if gt $i 0}},{{end}}{
                        "address": "{{.server}}",
                        "port": {{.port}},
                        "users": [
                            {
                                "id": "{{.uuid}}"
                            }
                        ]
                    }{{ $i = 1 }}{{end}}{{end}}
                ]
            }
        },{
            "protocol": "freedom",
            "tag": "direct"
        }
    ],
    "routing": {
        "domainStrategy": "IPOnDemand",
        "rules": [
            {
                "type": "field",
                "ip": [
                    "geoip:private"
                ],
                "outboundTag": "direct"
            }
        ]
    }
}