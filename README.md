[![CC BY-NC-SA 4.0][cc-by-nc-sa-shield]][cc-by-nc-sa]

本作品采用
[知识共享署名-非商业性使用-相同方式共享 4.0 国际许可协议][cc-by-nc-sa]进行许可。

[cc-by-nc-sa]: http://creativecommons.org/licenses/by-nc-sa/4.0/
[cc-by-nc-sa-shield]: https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg

## 后端 Docker 部署指南

### 使用镜像
```bash
docker pull swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:latest
```

### 运行容器
```bash
docker run -d \
  --name keyframe-back \
  -v /chaozj/data/etc:/etc:rw \
  -p 6470:6470 \
  -p 8888:8888 \
  swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:latest
```

### 端口说明
- 8888: 后端服务端口
- 6470: Prometheus监控指标服务端口

### 目录挂载说明
- 将主机目录`/chaozj/data/etc`挂载到容器内的`/etc`目录，权限为读写(rw)
- 该目录用于存储配置文件`keyframeGo.yaml`和IP数据库`cz88_public_v4.czdb`和`cz88_public_v6.czdb`
#### 配置文件
- 容器启动时会自动加载`/etc/keyframeGo.yaml`配置文件
- 请确保配置文件中的敏感信息已正确配置

## 后端 二进制文件部署指南
### 下载二进制文件
根据您的系统架构选择对应的二进制文件：
- macOS (ARM64): [keyframeGo-darwin-arm64](https://github.com/isTrih/KeyFrame/releases/latest)
- Linux (ARM64): [keyframeGo-linux-arm64](https://github.com/isTrih/KeyFrame/releases/latest)
- Linux (AMD64): [keyframeGo-linux-amd64](https://github.com/isTrih/KeyFrame/releases/latest)

### 部署步骤
1. 将下载的二进制文件移动到您的工作目录
2. 在工作目录下创建etc目录：`mkdir -p etc`
3. 在etc目录下创建keyframeGo.yaml配置文件
4. 下载IP数据库文件并放入etc目录：
   - [cz88_public_v4.czdb](https://www.cz88.net/geo-public)
   - [cz88_public_v6.czdb](https://www.cz88.net/geo-public)
5. 授予二进制文件执行权限：`chmod +x keyframeGo-*`
6. 运行程序：`./keyframeGo-*`

注意：请确保配置文件中的路径与您实际部署的目录结构一致。

## keyframeGo.yaml说明
```yaml
# /etc/keyframeGo.yaml
# 基础配置
Name: keyframeGo  # 服务名称
Host: 0.0.0.0     # 监听地址
Port: 8888        # 服务端口

# 开发服务器配置
DevServer:
  Enabled: true   # 是否启用开发模式
  Port: 6470      # Prometheus监控端口
  MetricsPath: /metrics  # 监控指标路径
  EnableMetrics: true    # 是否启用监控

# Redis缓存配置
Cache:
  - Host: 127.0.0.1:6379  # Redis地址
    Pass: your_password   # Redis密码

# MySQL数据库配置
DB:
  DataSource: username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# PostgreSQL配置
PG:
  DataSource: postgresql://username:password@host:port/database

# JWT认证配置
Auth:
  AccessSecret: your_jwt_secret  # JWT签名密钥
  AccessExpire: 604800           # 过期时间(秒)

# 七牛云存储配置
Qiniu:
  AK: your_access_key     # 七牛AccessKey
  SK: your_secret_key     # 七牛SecretKey

# 短信服务配置
Unisms:
  SK: your_sms_secret_key  # 短信服务密钥

# IP检查配置
IPCheck:
  Path4: etc/cz88_public_v4.czdb  # IPv4数据库路径
  Path6: etc/cz88_public_v6.czdb  # IPv6数据库路径
  KEY: your_ipcheck_key           # 解密密钥

# NATS消息队列配置
NATS:
  ADDR: host:port  # NATS服务器地址

# 注意：实际配置时请替换示例中的占位符，并确保敏感信息的安全性
```
