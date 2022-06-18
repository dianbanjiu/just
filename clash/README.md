此处是一个基础的 clash 的配置文件，你可以直接将其复制为你的 clash 配置文件，并修改其中的 `proxies` 和 `proxy-group`  为你自己的服务配置即可

clash 默认的配置文件位置在 `~/.config/clash/config.yaml`  

修改完成之后记得使用 clash 程序检查一下配置文件的配置是否正确
```shell
clash -t -f ~/.config/clash/config.yaml
```

此处提供的配置文件看起来可能有点混乱，是因为这个文件已经使用此程序更新过订阅信息，写入的过程中组件根据键名对文件进行了一次排序
