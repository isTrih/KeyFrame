# 使用 scratch 作为基础镜像
FROM swr.cn-east-3.myhuaweicloud.com/keyframe/alpine:latest
LABEL authors="trih"

# 设置时区为 Asia/Shanghai
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 将你的二进制文件复制到容器中
COPY prometheus-nats-exporter /prometheus-nats-exporter
COPY /etc/chaozjani.yaml /etc/keyframeGo.yaml
COPY /etc/cz88_public_v4.czdb etc/cz88_public_v4.czdb
COPY /etc/cz88_public_v6.czdb etc/cz88_public_v6.czdb
# 设置可执行权限（如果需要）
RUN chmod +x /keyframeGo

# 指定容器启动时运行的命令
CMD ["/keyframeGo"]