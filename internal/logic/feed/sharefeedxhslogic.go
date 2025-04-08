package feed

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/x/errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareFeedXHSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewShareFeedXHSLogic // 分享小红书帧（文章）
func NewShareFeedXHSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareFeedXHSLogic {
	return &ShareFeedXHSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareFeedXHSLogic) ShareFeedXHS() (resp *types.ShareFeedXHSResponse, err error) {
	appKey := "red.jdTVXR4Ldj9sudhb"
	appSecret := "e01fca2788b5d8097d76fbec72615e00"

	// 获取 access_token
	accessToken, expiresIn, err := getAccessToken(appKey, appSecret)
	if err != nil {
		fmt.Printf("获取 access_token 失败: %v\n", err)
		return nil, errors.New(4003, "获取 access_token 失败")
	}

	fmt.Printf("access_token: %s\n", accessToken)
	fmt.Printf("expires_in: %d\n", expiresIn)

	// 生成第二次签名示例
	nonce := generateNonce()
	timestamp := time.Now().UnixMilli()
	timestampStr := strconv.FormatInt(timestamp, 10)
	signature := generateSignature(appKey, nonce, timestampStr, accessToken)
	return &types.ShareFeedXHSResponse{
		Nonce:     nonce,
		Timestamp: timestampStr,
		Signature: signature,
	}, nil
}

// 生成随机字符串
func generateNonce() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nonce := make([]byte, 32)
	for i := range nonce {
		nonce[i] = charset[rand.Intn(len(charset))]
	}
	return string(nonce)
}

// 生成签名
func generateSignature(appKey, nonce, timestamp, appSecret string) string {
	// 拼接参数
	paramStr := fmt.Sprintf("appKey=%s&nonce=%s&timeStamp=%s%s", appKey, nonce, timestamp, appSecret)
	// 计算 SHA-256 哈希
	hash := sha256.Sum256([]byte(paramStr))
	return hex.EncodeToString(hash[:])
}

// 获取 access_token
func getAccessToken(appKey, appSecret string) (string, int64, error) {
	// 生成随机字符串和时间戳
	nonce := generateNonce()
	timestamp := time.Now().UnixMilli()
	timestampStr := strconv.FormatInt(timestamp, 10)

	fmt.Println(nonce)
	fmt.Println(timestamp)
	fmt.Println(timestampStr)

	// 生成签名
	signature := generateSignature(appKey, nonce, timestampStr, appSecret)
	fmt.Println(signature)
	// 构建请求参数
	requestBody := map[string]interface{}{
		"app_key":   appKey,
		"nonce":     nonce,
		"timestamp": timestampStr,
		"signature": signature,
	}

	// 将请求参数转换为 JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, err
	}

	// 创建一个自定义的 HTTP 客户端，忽略 SSL 证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// 发送 POST 请求
	resp, err := client.Post("https://edith.xiaohongshu.com/api/sns/v1/ext/access/token", "application/json",
		bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("关闭响应体失败: %v\n", err)
		}
	}(resp.Body)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	// 定义响应结构体
	type Response struct {
		Code    int    `json:"code"`
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Data    struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int64  `json:"expires_in"`
		} `json:"data"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", 0, err
	}

	if !response.Success {
		return "", 0, fmt.Errorf("获取 access_token 失败，错误码: %d，错误信息: %s", response.Code, response.Msg)
	}

	return response.Data.AccessToken, response.Data.ExpiresIn, nil
}
