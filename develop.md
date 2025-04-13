[//]: # (- application 目录存放微服务)

[//]: # (- db 目录存放数据库文件)

[//]: # (  - applet http接口)

[//]: # (  - article 文章rpc)

[//]: # (  - user 用户rpc)

[//]: # (  - concerned 关注rpc)

[//]: # (  - member 会员rpc)

[//]: # (  - message 消息rpc)

[//]: # (- pkg 目录存放通用方法)

需要将chaozjani.yaml文件拷贝到项目根目录下


```
用户请求 → Go-zero API 服务 → Redis (HINCRBY) → NATS → Go-zero Consumer → MySQL (批量更新)
│                                      │
└─ 兜底直接更新 MySQL (fallbackToDB)     └─ 定时一致性校验
```


```go
//事务使用方法
    var conn sqlx.SqlConn
    err := conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        r, err := session.ExecCtx(ctx, "insert into user (id, name) values (?, ?)", 1, "test")
        if err != nil {
            return err
        }
        r ,err =session.ExecCtx(ctx, "insert into user (id, name) values (?, ?)", 2, "test")
        if err != nil {
            return err
        }
    })
```