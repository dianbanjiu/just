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
            "tag": "proxy",
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
    "dns": {
    "servers": [
        {
            "address": "1.1.1.1",
            "domains": ["geosite:geolocation-!cn"]
        },
        {
            "address": "223.5.5.5",
            "domains": ["geosite:cn"],
            "expectIPs": ["geoip:cn"]
        },
        {
            "address": "114.114.114.114",
            "domains": ["geosite:cn"]
        },
        "localhost"
        ]
    },
    "routing": {
        "domainStrategy": "IPOnDemand",
        "rules": [
            {
                "type": "field",
                "ip": [
                    "geoip:cn",
                    "geoip:private"
                ],
                "outboundTag": "direct"
            },{
                "type": "field",
                "domain": ["geosite:geolocation-!cn"],
                "outboundTag": "proxy"
            },{
                "type": "field",
                "ip": ["223.5.5.5"],
                "outboundTag": "direct"
            }
        ]
    }
}