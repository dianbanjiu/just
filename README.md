# just
## 说明
just 是一个 justmysocks 的实用工具，可以使用 justmysocks 提供的 API 来解析订阅信息、打印流量使用情况  

此项目是基于我自己的订阅解析内容实现的，不一定适用于所有人

so  
**慎用写入功能**  
**慎用写入功能**  
**慎用写入功能**  
OR  
**写入之前记得备份**  
**写入之前记得备份**  
**写入之前记得备份**  


just 只会请求 justmysocks 的提供的订阅地址和 API，不会向任何第三方发送任何请求。如果 justmysocks 的订阅地址或者 API 被墙了的话，此程序可能会运行失败 (～￣▽￣)～

## 使用方式
复制仓库下面的 config.sample.yaml 为 config.yaml，从 justmysocks 账户中将自己的订阅地址和流量使用 API 粘贴到 config.yaml 中

复制 [examples](/examples) 下需要的文件到项目根目录下并重命名为 proxy.tpl，根据自己的需要修改 proxy.tpl  

**打印流量使用情况**
```shell
just usage
```

**将订阅信息写入到指定位置**，这里写入时支持两种方式进行配置
1. 修改 just 的配置文件，将 `proxy_config_path` 修改为自己配置文件的位置
2. 通过 `-c` 参数来指定配置文件的位置


```shell
just sub -w -c ~/path/to/config.yaml
```
