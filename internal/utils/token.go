package utils

import "github.com/golang-jwt/jwt/v4"

// GetToken 生成token, iat为unix时间戳, secretKey为密钥, payloads为负载, seconds为过期时间(秒)
// 示例:
//
//	 payloads := make(map[string]any)
//		payloads["UID"], _ = user.LastInsertId()
//		payloads["UTYPE"] = 0
//	 payloads不加密，不能放敏感信息
func GetToken(iat int64, secretKey string, payloads map[string]any, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["expTime"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
