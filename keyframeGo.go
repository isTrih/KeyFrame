package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"zerobackend/internal/config"
	"zerobackend/internal/handler"
	"zerobackend/internal/svc"
)

var configFile = flag.String("f", "etc/keyframeGo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 设置允许跨域的域名
	//# 需要通过的域名，这里可以写多个域名 或者可以写 * 全部通过
	//"http://127.0.0.1", "https://go-zero.dev", "http://localhost", "*.chaozj.com", "https://ani.chaozj.com"}
	domains := []string{"*"}
	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors(domains...),
		rest.WithCustomCors(func(header http.Header) {
			header.Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id,OS,Platform, Version,kip")
			header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
			header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, ")
		}, nil, "*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//// 初始化 NATS 连接
	//natsclient.InitNats()
	//// 启动消息消费者
	//go consumer.StartMessageConsumer()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	defer ctx.IP4Searcher.Close()
	server.Start()
}
