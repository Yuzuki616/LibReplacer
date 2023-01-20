# LibReplacer

一个用于替换Emby库中内容的反代服务器。
本项目支持自动申请tls证书及Https，可在某种程度上替代nginx反代。

# 使用

从Releases里下载可执行文件，修改默认配置文件后直接执行或使用systemd等启动。
配置文件中`LibraryPath`项填写的路径会被监听变动，被修改后会自动重载，可利用外部程序生成榜单数据输出至此路径，具体格式请查看example。


# 构建

``` bash
git clone https://github.com/Yuzuki616/LibReplacer.git
cd LibReplacer
go build
```

# Thanks

- [Lego](https://github.com/go-acme/lego)
- [Gin](https://github.com/gin-gonic/gin)
