# just
## 说明
just 是一个 justmysocks 的实用工具，可以使用 justmysocks 提供的 API 来解析订阅信息、打印流量使用情况  

## 使用方式
复制仓库下面的 config.sample.yaml 为 config.yaml，从 justmysocks 账户中将自己的订阅地址和流量使用 API 粘贴到 config.yaml 中

打印流量使用情况
```shell
just usage
```

打印订阅信息
```shell
just sub -p
```

将订阅信息写入到 clash 的配置文件中去，这里写入时支持两种方式
1. 修改 just 的配置文件，将 `clash_config_path` 修改为自己 clash 配置文件的位置
2. 通过 `-c` 参数来指定 clash 配置文件的位置


为了避免意外，可以先将 clash 的配置文件先做一个备份，然后再进行写入操作  
```shell
just sub -w -c ~/.config/clash/config.yaml
```

检查配置文件是否有误
```shell
clash -t -f ~/.config/clash/config.yaml
```
