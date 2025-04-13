package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"net/http"
	"strconv"
	"zerobackend/internal/config"
)

// Message 定义消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// InspReq 定义请求结构体
type InspReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Key      string    `header:"Authorization"`
}

// Choice 定义响应中 choices 数组里的元素结构体
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 定义响应中 usage 结构体
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// InspResp 定义完整的响应结构体
type InspResp struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

func DoInsp(c config.Config, Content string) (uint8, error) {
	domain := c.Insp.URL
	model := c.Insp.MODEL

	// 构建请求体
	req := InspReq{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是文本审核员，判断信息是否出现1:色情, 2:骚扰, 3:广告, 4:政治, 5:引战, 6:辱骂, 0:正常，并仅返回一个数字。（请注意谐音，emoji等规避方法，带有流量卡和其他引人不适的广告内容即视为1广告）",
			}, {
				Role:    "user",
				Content: Content,
			}},
		Key: "Bearer " + c.Insp.KEY,
	}

	// 创建自定义的 http.Transport 并忽略证书验证
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// 创建自定义的 http.Client
	customClient := &http.Client{
		Transport: transport,
	}

	// 使用自定义的 http.Client 创建服务
	serviceName := "ai_insp"
	service := httpc.NewServiceWithClient(serviceName, customClient)
	// 发送请求

	resp, err := service.Do(context.Background(), http.MethodPost, domain, req)
	if err != nil {
		fmt.Println("出错：", err)
		return 99, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body) // 确保在函数结束时关闭响应体

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体出错：", err)
		return 99, err
	}

	// 反序列化响应内容到 InspResp 结构体
	var inspResp InspResp
	err = json.Unmarshal(body, &inspResp)
	if err != nil {
		fmt.Println("反序列化响应内容出错：", err)
		return 99, err
	}

	// 检查 choices 数组是否为空
	if len(inspResp.Choices) == 0 {
		return 99, fmt.Errorf("响应中 choices 数组为空")
	}

	// 获取 message 里的 content 字符串
	contentStr := inspResp.Choices[0].Message.Content

	// 将字符串转换为 uint8
	result, err := strconv.ParseUint(contentStr, 10, 8)
	if err != nil {
		fmt.Println("将 content 转换为 uint8 出错：", err)
		return 99, err
	}

	return uint8(result), nil
}
