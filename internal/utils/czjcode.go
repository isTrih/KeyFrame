package utils

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DecryptTriDESToNumber 解密Triple DES加密的字符串并转换为数字
// key: 密钥
// iv: 初始向量
// encryptedHex: 加密的十六进制字符串(如"4bb843c91fa23ef8")
// 成功解密返回数字，失败返回0
func DecryptTriDESToNumber(key, iv, encryptedHex string) int64 {
	// 输入参数校验
	if len(key) == 0 || len(iv) == 0 {
		fmt.Printf("DecryptTriDESToNumber error: empty key or iv\n")
		return 0
	}

	if len(encryptedHex) == 0 {
		fmt.Printf("DecryptTriDESToNumber error: empty encrypted hex\n")
		return 0
	}

	// 解码十六进制字符串
	encryptedData, err := hex.DecodeString(encryptedHex)
	if err != nil {
		fmt.Printf("DecryptTriDESToNumber hex.DecodeString error: %v\n", err)
		return 0
	}

	// 检查加密数据长度是否合法(需要是8的倍数)
	if len(encryptedData) == 0 || len(encryptedData)%8 != 0 {
		fmt.Printf("DecryptTriDESToNumber error: invalid encrypted data length: %d\n", len(encryptedData))
		return 0
	}

	// 填充密钥到24字节(Triple DES需要24字节密钥)
	paddedKey := padKeyTo24Bytes([]byte(key))

	// 填充IV到8字节(DES块大小为8字节)
	paddedIV := padIVTo8Bytes([]byte(iv))

	// 创建Triple DES密码块
	block, err := des.NewTripleDESCipher(paddedKey)
	if err != nil {
		fmt.Printf("DecryptTriDESToNumber NewTripleDESCipher error: %v\n", err)
		return 0
	}

	// 使用CBC模式解密
	mode := cipher.NewCBCDecrypter(block, paddedIV)
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// 去除PKCS#7填充
	decrypted, err = removePKCS7Padding(decrypted)
	if err != nil {
		fmt.Printf("DecryptTriDESToNumber removePKCS7Padding error: %v\n", err)
		return 0
	}

	// 转换为字符串并清理空白字符
	decryptedStr := strings.TrimSpace(string(decrypted))
	if len(decryptedStr) == 0 {
		fmt.Printf("DecryptTriDESToNumber error: empty decrypted string\n")
		return 0
	}

	// 尝试将解密结果转换为数字
	num, err := strconv.ParseInt(decryptedStr, 10, 64)
	if err != nil {
		fmt.Printf("DecryptTriDESToNumber ParseInt error: %v, decrypted string: '%s'\n", err, decryptedStr)
		return 0
	}

	return num
}

// padKeyTo24Bytes 将密钥填充到24字节(Triple DES需要)
func padKeyTo24Bytes(key []byte) []byte {
	if len(key) >= 24 {
		return key[:24]
	}

	result := make([]byte, 24)
	copy(result, key)

	// 重复密钥填满24字节
	for i := len(key); i < 24; i++ {
		result[i] = key[i%len(key)]
	}

	return result
}

// padIVTo8Bytes 将IV填充到8字节
func padIVTo8Bytes(iv []byte) []byte {
	if len(iv) >= 8 {
		return iv[:8]
	}

	result := make([]byte, 8)
	copy(result, iv)

	// 重复IV填满8字节
	for i := len(iv); i < 8; i++ {
		result[i] = iv[i%len(iv)]
	}

	return result
}

// removePKCS7Padding 移除PKCS#7填充
func removePKCS7Padding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}

	paddingLen := int(data[length-1])
	if paddingLen > 8 || paddingLen == 0 {
		return nil, fmt.Errorf("invalid padding length: %d", paddingLen)
	}

	// 验证填充
	for i := length - paddingLen; i < length; i++ {
		if data[i] != byte(paddingLen) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:length-paddingLen], nil
}
