# Port of HTTP(S) proxy server on the local end
port: 1007

# Port of SOCKS5 proxy server on the local end
socks-port: 1006

# Transparent proxy server port for Linux and MacOS(Redirect TCP and TProxy UDP)
# redir-port: 7892

# Transparent proxy server port for Linux(TProxy TCP and UDP)
# tproxy-port: 7893

# HTTP(S) and SOCKS4(A)/SOCKS5 server on the same port
# mixed-port: 7890

# authentication of local SOCKS5/HTTP(S) server
# authentication:
#   - "user1: pass1"
#   - "user2: pass2"

# set to true to allow connections to the local-end server form other LAN IP address
allow-lan: false

# This is only applicable when `allow-lan` is `true`
# '*': bind all IP address
# 192.168.199.199: bind a single IPV4 address
# '[aaaa::a8aa:ff:fe09:57d8]': bind a single IPV6 address
bind-address: '*'

# clash routerworking mode 
# rule: rule-based packet routing
# global: all packets will forwarded to a single endpoint
# direct: directly forward the packets to the Internet
mode: rule

# clash by defult prints logs to the STDOUT
# info / warning / error / debug / silent
log-level: info

# When set to false, resolver won't translate hostnames to IPv6 addresses
ipv6: false

# RESTful web API listening address
# external-controller: 127.0.0.1:9090

# Secret for the RESTful API (optional)
# Authenticate by spedifying HTTP header `Authorization: Bearer ${secret}`
# ALWAYS set a secret if RESTful API is listening on 0.0.0.0
# secret: ""

# Outbound interface name
# interface-name: en0

# fwmark on Linux only
# routing-mark: 6666

# Static hosts for DNS server and connection establishment (like /etc/hosts)
#
# Wildcard hostnames are supported (e.g. *.clash.dev, *.foo.*.example.com)
# Non-wildcard domain names have a higher priority than wildcard domain names
# e.g. foo.example.com > *.example.com > .example.com
# P.S. +.foo.com equals to .foo.com and foo.com
hosts:
  # '*.clash.dev': 127.0.0.1
  # '.dev': 127.0.0.1
  # 'alpha.clash.dev': '::1'

profile:
  # Store the `select` results in $HOME/.config/clash/.cache
  # set false If you don't want this behavior
  # when two different configurations have groups with the same name, the selected values are shared
  store-selected: false

  # persistence fakeip
  store-fake-ip: true

# DNS server settings
# This section is optional. When not present, the DNS server will be disabled.
dns:
  enable: false
  listen: 0.0.0.0:53
  # ipv6: false # when the false, response to AAAA questions will be empty

  # These nameservers are used to resolve the DNS nameserver hostnames below.
  # Specify IP addresses only
  default-nameserver:
    - 114.114.114.114
    - 8.8.8.8
  enhanced-mode: fake-ip # or redir-host (not recommended)
  fake-ip-range: 198.18.0.1/16 # Fake IP addresses pool CIDR
  # use-hosts: true # lookup hosts and return IP record
  
  # Hostnames in this list will not be resolved with fake IPs
  # i.e. questions to these domain names will always be answered with their
  # real IP addresses
  # fake-ip-filter:
  #   - '*.lan'
  #   - localhost.ptlogin2.qq.com
  
  # Supports UDP, TCP, DoT, DoH. You can specify the port to connect to.
  # All DNS questions are sent directly to the nameserver, without proxies
  # involved. Clash answers the DNS question with the first result gathered.
  nameserver:
    - 114.114.114.114 # default value
    - 8.8.8.8 # default value
    - tls://dns.rubyfish.cn:853 # DNS over TLS
    - https://1.1.1.1/dns-query # DNS over HTTPS
    - dhcp://en0 # dns from dhcp
    # - '8.8.8.8#en0'

  # When `fallback` is present, the DNS server will send concurrent requests
  # to the servers in this section along with servers in `nameservers`.
  # The answers from fallback servers are used when the GEOIP country
  # is not `CN`.
  # fallback:
  #   - tcp://1.1.1.1
  #   - 'tcp://1.1.1.1#en0'

  # If IP addresses resolved with servers in `nameservers` are in the specified
  # subnets below, they are considered invalid and results from `fallback`
  # servers are used instead.
  #
  # IP address resolved with servers in `nameserver` is used when
  # `fallback-filter.geoip` is true and when GEOIP of the IP address is `CN`.
  #
  # If `fallback-filter.geoip` is false, results from `nameserver` nameservers
  # are always used if not match `fallback-filter.ipcidr`.
  #
  # This is a countermeasure against DNS pollution attacks.
  # fallback-filter:
  #   geoip: true
  #   geoip-code: CN
  #   ipcidr:
  #     - 240.0.0.0/4
  #   domain:
  #     - '+.google.com'
  #     - '+.facebook.com'
  #     - '+.youtube.com'
  
  # Lookup domains via specific nameservers
  # nameserver-policy:
  #   'www.baidu.com': '114.114.114.114'
  #   '+.internal.crop.com': '10.0.0.1'

proxies:
{{range .}}
  - name: {{.name}}
    type: {{.type}}
    server: {{.server}}
    port: {{.port}}
    cipher: {{.cipher}}
    {{if .password}}password: {{.password}}{{end}}
    {{if .uuid}}uuid: {{.uuid}}{{end}}
    {{if .alterId}}alterId: {{.alterId}}{{end}}
{{end}}
proxy-groups:
  # url-test select which proxy will be used by benchmarking speed to a URL.
  - name: "auto"
    type: url-test
    proxies:
    {{range .}}
      - {{.name}}{{end}}
    # tolerance: 150
    # lazy: true
    url: 'http://www.gstatic.com/generate_204'
    interval: 300

  # fallback selects an available policy by priority. The availability is tested by accessing an URL, just like an auto url-test group.
  - name: "fallback-auto"
    type: fallback
    proxies:
    {{range .}}
      - {{.name}}{{end}}
    url: 'http://www.gstatic.com/generate_204'
    interval: 300

  # load-balance: The request of the same eTLD+1 will be dial to the same proxy.
  - name: "load-balance"
    type: load-balance
    proxies:
    {{range .}}
      - {{.name}}{{end}}
    url: 'http://www.gstatic.com/generate_204'
    interval: 300
    # strategy: consistent-hashing # or round-robin

  # select is used for selecting proxy or proxy group
  # you can use RESTful API to switch proxy is recommended for use in GUI.
  - name: Proxy
    type: select
    # disable-udp: true
    proxies:
      - auto

rules:
  - DOMAIN-KEYWORD,amazon,Proxy
  - DOMAIN-KEYWORD,google,Proxy
  - DOMAIN-KEYWORD,gmail,Proxy
  - DOMAIN-KEYWORD,youtube,Proxy
  - DOMAIN-KEYWORD,facebook,Proxy
  - DOMAIN-SUFFIX,fb.me,Proxy
  - DOMAIN-SUFFIX,fbcdn.net,Proxy
  - DOMAIN-KEYWORD,twitter,Proxy
  - DOMAIN-KEYWORD,instagram,Proxy
  - DOMAIN-KEYWORD,dropbox,Proxy
  - DOMAIN-SUFFIX,twimg.com,Proxy
  - DOMAIN-KEYWORD,blogspot,Proxy
  - DOMAIN-SUFFIX,youtu.be,Proxy
  - DOMAIN-KEYWORD,whatsapp,Proxy
  - DOMAIN-KEYWORD,admarvel,REJECT
  - DOMAIN-KEYWORD,admaster,REJECT
  - DOMAIN-KEYWORD,adsage,REJECT
  - DOMAIN-KEYWORD,adsmogo,REJECT
  - DOMAIN-KEYWORD,adsrvmedia,REJECT
  - DOMAIN-KEYWORD,adwords,REJECT
  - DOMAIN-KEYWORD,adservice,REJECT
  - DOMAIN-KEYWORD,domob,REJECT
  - DOMAIN-KEYWORD,duomeng,REJECT
  - DOMAIN-KEYWORD,dwtrack,REJECT
  - DOMAIN-KEYWORD,guanggao,REJECT
  - DOMAIN-KEYWORD,lianmeng,REJECT
  - DOMAIN-SUFFIX,mmstat.com,REJECT
  - DOMAIN-KEYWORD,omgmta,REJECT
  - DOMAIN-KEYWORD,openx,REJECT
  - DOMAIN-KEYWORD,partnerad,REJECT
  - DOMAIN-KEYWORD,pingfore,REJECT
  - DOMAIN-KEYWORD,supersonicads,REJECT
  - DOMAIN-KEYWORD,uedas,REJECT
  - DOMAIN-KEYWORD,umeng,REJECT
  - DOMAIN-KEYWORD,usage,REJECT
  - DOMAIN-KEYWORD,wlmonitor,REJECT
  - DOMAIN-KEYWORD,zjtoolbar,REJECT
  - DOMAIN-SUFFIX,local,DIRECT
  - IP-CIDR,127.0.0.0/8,DIRECT
  - IP-CIDR,172.16.0.0/12,DIRECT
  - IP-CIDR,192.168.0.0/16,DIRECT
  - IP-CIDR,10.0.0.0/8,DIRECT
  - IP-CIDR,17.0.0.0/8,DIRECT
  - IP-CIDR,100.64.0.0/10,DIRECT
  - GEOIP,CN,DIRECT
  - MATCH,Proxy