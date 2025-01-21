# 使用 scratch 作为基础镜像
FROM alpine:latest
LABEL authors="trih"

# 设置时区为 Asia/Shanghai
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 将你的二进制文件复制到容器中
COPY chaozjani /chaozjani
COPY etc /etc
# 设置可执行权限（如果需要）
RUN chmod +x /chaozjani

# 指定容器启动时运行的命令
CMD ["/chaozjani"]
