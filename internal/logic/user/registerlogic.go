package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"time"
	"zerobackend/internal/utils"
	model "zerobackend/mdl/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/x/errors"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	insp, err := utils.DoInsp(l.svcCtx.Config, req.Username)
	if insp != 0 {
		// 违规逻辑
		return nil, errors.New(6099, "昵称有违规内容")
	}

	rds, rds2 := utils.RedisCheck(l.svcCtx.Config, req.Mobile, req.VerifyCode)
	if rds != nil {
		return nil, rds
	}
	if rds2 != nil {
		return nil, rds2
	}

	focusId := utils.DecryptTriDESToNumber(l.svcCtx.Config.InviteKey.KEY,
		l.svcCtx.Config.InviteKey.IV, req.CZJCode)
	fmt.Println("获取邀请码", focusId)
	if focusId != 0 {
		IDcheck, IDcheckerr := l.svcCtx.UserModel.FindOne(l.ctx, focusId)
		if IDcheckerr != nil && IDcheckerr != model.ErrNotFound {
			return nil, errors.New(4003, "查询数据失败")
		}
		if IDcheck == nil {
			check, checkerr := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
			if checkerr != nil && checkerr != model.ErrNotFound {
				return nil, errors.New(4003, "查询数据失败")
			}
			if check == nil {
				err = l.svcCtx.UserModel.CreateUserByCZJCode(
					l.ctx, focusId,
					req.Username, EncryptPassword(req.Password),
					req.Mobile)
				if err != nil {
					fmt.Println(err)
					return nil, errors.New(6011, "注册失败")
				}

				payloads := make(map[string]any)

				// 现在可以安全访问 uid.Id
				payloads["UID"] = focusId
				payloads["UTYPE"] = 0
				payloads["USTATUS"] = 0

				accessToken, tokenErr := utils.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
				if tokenErr != nil {
					return nil, tokenErr
				}

				// 改进 IP 处理部分
				var ip = req.KIP
				if ip == "" {
					// 跳过 IP 处理
				} else {
					parsedIP := net.ParseIP(ip)
					if parsedIP == nil {
						// 跳过 IP 处理
					} else if parsedIP.To4() != nil {
						region, err := l.svcCtx.IP4Searcher.Search(ip)
						if err != nil {
							// 记录错误但不中断流程
						} else {
							// 更新 IP 信息
							upErr := l.svcCtx.UserModel.UpdateIpByMobile(l.ctx, req.Mobile, region, ip)
							if upErr != nil {
								// 记录错误但不中断流程
							}
						}
					} else if parsedIP.To16() != nil {
						region, err := l.svcCtx.IP6Searcher.Search(ip)
						if err != nil {
							// 记录错误但不中断流程
						} else {
							// 更新 IP 信息
							upErr := l.svcCtx.UserModel.UpdateIpByMobile(l.ctx, req.Mobile, region, ip)
							if upErr != nil {
								// 记录错误但不中断流程
							}
						}
					}
				}

				resp = new(types.RegisterResponse)
				resp.UserId = uint64(focusId)
				resp.UserName = req.Username
				resp.Avatar = "avatar.jpg"
				resp.Signature = ""
				resp.Token = accessToken
				resp.UserType = 0
				err = utils.SMSDelete(l.svcCtx.Config, req.Mobile)
				if err != nil {
					return nil, err
				}
				return resp, nil
			} else {
				return nil, errors.New(6012, "用户已存在")
			}
		} else {
			return nil, errors.New(6012, "邀请码已使用")
		}
	} else {
		check, checkerr := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
		if checkerr != nil && checkerr != model.ErrNotFound {
			return nil, errors.New(4003, "查询数据失败")
		}
		if check == nil {
			_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
				Mobile:    req.Mobile,
				Nickname:  req.Username,
				Password:  EncryptPassword(req.Password),
				Signature: "",
				Avatar:    "avatar.jpg",
			})
			if err != nil {
				fmt.Println(err)
				return nil, errors.New(6011, "注册失败")
			}

			payloads := make(map[string]any)

			// 2. 修改 FindOneByMobile 后的处理逻辑
			uid, findErr := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
			if findErr != nil {
				logc.Errorf(l.ctx, "查询用户信息失败: %v", findErr)
				return nil, errors.New(5001, "用户注册成功但查询失败")
			}

			if uid == nil {
				logc.Error(l.ctx, "用户查询结果为空")
				return nil, errors.New(5002, "用户信息异常")
			}

			// 现在可以安全访问 uid.Id
			uintUid := uint64(uid.Id)
			payloads["UID"] = uintUid
			payloads["UTYPE"] = uid.Type
			payloads["USTATUS"] = uid.Status

			accessToken, tokenErr := utils.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
			if tokenErr != nil {
				return nil, tokenErr
			}

			// 改进 IP 处理部分
			var ip = req.KIP
			if ip == "" {
				// 跳过 IP 处理
			} else {
				parsedIP := net.ParseIP(ip)
				if parsedIP == nil {
					// 跳过 IP 处理
				} else if parsedIP.To4() != nil {
					region, err := l.svcCtx.IP4Searcher.Search(ip)
					if err != nil {
						// 记录错误但不中断流程
					} else {
						// 更新 IP 信息
						upErr := l.svcCtx.UserModel.UpdateIpByMobile(l.ctx, req.Mobile, region, ip)
						if upErr != nil {
							// 记录错误但不中断流程
						}
					}
				} else if parsedIP.To16() != nil {
					region, err := l.svcCtx.IP6Searcher.Search(ip)
					if err != nil {
						// 记录错误但不中断流程
					} else {
						// 更新 IP 信息
						upErr := l.svcCtx.UserModel.UpdateIpByMobile(l.ctx, req.Mobile, region, ip)
						if upErr != nil {
							// 记录错误但不中断流程
						}
					}
				}
			}

			resp = new(types.RegisterResponse)
			resp.UserId = uintUid
			resp.UserName = req.Username
			resp.Avatar = "avatar.jpg"
			resp.Signature = uid.Signature
			resp.Token = accessToken
			resp.UserType = uint16(uid.Type)
			err = utils.SMSDelete(l.svcCtx.Config, req.Mobile)
			if err != nil {
				return nil, err
			}
			return resp, nil
		} else {
			return nil, errors.New(6012, "用户已存在")
		}
	}
}

// EncryptPassword 加密密码
// 示例:
//
//	encryptedPassword := EncryptPassword("123456")
func EncryptPassword(needEncryptPassword string) (encryptedPassword string) {

	encryptedPassword = fmt.Sprintf("%x", sha256.Sum256([]byte(needEncryptPassword+"tmh")))
	return encryptedPassword
}
