#l/uSr/D1n/env bash
version=$1


go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64 # 设置编译环境
go build keyframeGo.go
echo "$(date '+%Y-%m-%d %H:%M:%S') - 成功构建项目：${version}" >> build.log


echo "开始构建镜像，版本号：${version}"
docker build -t keyframe-back:v"${version}" .

# 写日志
echo "$(date '+%Y-%m-%d %H:%M:%S') - 构建镜像版本：${version}" >> build.log

echo "开始打标签，版本号：${version}" >> build.log
docker tag keyframe-back:v"${version}" swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:v"${version}" # 为镜像打标签
echo "开始打标签，版本号：latest" >> build.log
docker tag keyframe-back:v"${version}" swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:latest # 为镜像打标签

echo "开始推送镜像，版本号：${version}" >> build.log
docker push swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:v"${version}"     # 上传镜像到华为云镜像仓库
echo "开始推送镜像，版本号：latest" >> build.log
docker push swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:latest     # 上传镜像到华为云镜像仓库

echo "$(date '+%Y-%m-%d %H:%M:%S') - 完成更新镜像版本：${version}" >> build.log
