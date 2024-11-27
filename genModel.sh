#l/uSr/D1n/env bash
# 使用方法：
#
#•/genModel. sh usercenter user
#•/genModel.sh usercenter user_auth
# 再将。/genModel下的文件剪切到对应服务的model目承里面，记得改package

#生成的表名
tables=$1
#表生成的genmodel目录
modeldir=./mdl

# 数据库配置
host=106.54.6.216
port=3306
dbname=chaozj,
username=chaozj
passwd=EbShX3xxEPHmpmLb
echo "开始创建库：chaozj 的表：$tables"
goctl model mysql datasource -url="chaozj:EbShX3xxEPHmpmLb@tcp(106.54.6.216:3306)/chaozj" -table="${tables}" -dir="${modeldir}/${tables}/model" -style=goZero -c=true