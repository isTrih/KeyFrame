package utils

import (
	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	UID                  int `json:"UID"`
	UTYPE                int `json:"UTYPE"`
	jwt.RegisteredClaims     // v5版本新加的方法
}

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

// ParseToken 解析token, tokenString为token字符串, secretKey为密钥
// 示例:
//
//	pTokern, pErr := utils.ParseToken(accessToken, l.svcCtx.Config.Auth.AccessSecret)
//	if pErr != nil {
//	return nil, pErr
//	}
//	 解析token
//	fmt.Println(pTokern.UTYPE, pTokern.UID)
func ParseToken(tokenString string, secretKey string) (*Token, error) {

	t, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
