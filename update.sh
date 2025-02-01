#l/uSr/D1n/env bash
version=$1

echo "开始构建镜像，版本号：${version}"
docker build -t keyframe-back:v"${version}" .

# 写日志
echo "$(date '+%Y-%m-%d %H:%M:%S') - 构建镜像版本：${version}" >> build.log

echo "开始推送镜像，版本号：${version}"
docker tag keyframe-back:v"${version}" swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:v"${version}" # 为镜像打标签
docker push swr.cn-east-3.myhuaweicloud.com/keyframe/keyframe-back:v"${version}"     # 上传镜像到华为云镜像仓库
echo "$(date '+%Y-%m-%d %H:%M:%S') - 推送镜像版本：${version}" >> build.log
