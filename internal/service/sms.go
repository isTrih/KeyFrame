package service

import (
	"fmt"
	unisms "github.com/apistd/uni-go-sdk/sms"
)

// SendSms 发送短信
func SendSms(mobile string, value string) (err error) {
	// 初始化
	client := unisms.NewClient("n6Z4ZVVUToSXF8ktQbXeNfGvQRp6mcdoH43pjNhy3uKPXMCW8") // 若使用简易验签模式仅传入第一个参数即可

	// 构建信息
	message := unisms.BuildMessage()
	message.SetTo(mobile)
	message.SetSignature("超正经科技")
	message.SetTemplateId("pub_verif_ttl3")
	message.SetTemplateData(map[string]string{"code": value, "ttl": "10"}) // 设置自定义参数 (变量短信)

	// 发送短信
	_, err = client.Send(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SendAllSms 发送短信（需要模版）
func SendAllSms(mobile string, value map[string]string, templateId string) (err error) {
	client := unisms.NewClient("n6Z4ZVVUToSXF8ktQbXeNfGvQRp6mcdoH43pjNhy3uKPXMCW8") // 若使用简易验签模式仅传入第一个参数即可

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
