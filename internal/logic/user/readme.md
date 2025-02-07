//TODO 完善

```go
uidjson, _ := l.ctx.Value("UID").(json.Number)
typejson, _ := l.ctx.Value("UTYPE").(json.Number)
statusjson, _ := l.ctx.Value("USTATUS").(json.Number)
uid, _ := uidjson.Int64()
utype, _ := typejson.Int64()
ustatus, _ := statusjson.Int64()

//user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
//if err != nil && err != model.ErrNotFound {
//	fmt.Println(err)
//	return nil, errors.New(4003, "查询数据失败")
//}
//if user == nil {
//	return nil, errors.New(6021, "用户不存在")
//}

resp = new(types.UserInfoResponse)
resp.UserId = uint64(uid)
resp.Username = "xx"
resp.Avatar = "XX"
resp.Type = uint8(utype)
resp.Status = uint8(ustatus)
return resp
```