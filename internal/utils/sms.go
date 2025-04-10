package utils

import (
	"fmt"
	"github.com/apistd/uni-go-sdk"
	unisms "github.com/apistd/uni-go-sdk/sms"
	"zerobackend/internal/config"
)

// SendSms 发送短信
func SendSms(c config.Config, mobile string, value string) (response *uni.UniResponse, err error) {
	// 初始化
	client := unisms.NewClient(c.Unisms.SK) // 若使用简易验签模式仅传入第一个参数即可

	// 构建信息
	message := unisms.BuildMessage()
	message.SetTo(mobile)
	message.SetSignature("芜湖超正经科技")
	message.SetTemplateId("pub_verif_ttl3")
	message.SetTemplateData(map[string]string{"code": value, "ttl": "10"}) // 设置自定义参数 (变量短信)

	// 发送短信
	res, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
}

// SendAllSms 发送短信（需要模版）
func SendAllSms(c config.Config, mobile string, value map[string]string, templateId string) (err error) {
	client := unisms.NewClient(c.Unisms.SK) // 若使用简易验签模式仅传入第一个参数即可

	// 构建信息
	message := unisms.BuildMessage()
	message.SetTo(mobile)
	message.SetSignature("超正经科技")
	message.SetTemplateId(templateId)
	message.SetTemplateData(value) // 设置自定义参数 (变量短信)

	// 发送短信
	_, err = client.Send(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
