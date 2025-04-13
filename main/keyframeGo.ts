import webapi from "./gocliRequest"
import * as components from "./keyframeGoComponents"
export * from "./keyframeGoComponents"

/**
 * @description "封禁动态"
 * @param req
 */
export function banFeed(req: components.ChangeUserTypeRequest) {
	return webapi.post<components.AdminRes>(`/v1/admin/feed/ban`, req)
}

/**
 * @description "解禁动态"
 * @param req
 */
export function unBanFeed(req: components.ChangeUserTypeRequest) {
	return webapi.post<components.AdminRes>(`/v1/admin/feed/unban`, req)
}

/**
 * @description "封禁用户"
 * @param req
 */
export function banUser(req: components.ChangeUserTypeRequest) {
	return webapi.post<components.AdminRes>(`/v1/admin/user/ban`, req)
}

/**
 * @description "改变用户权限"
 * @param req
 */
export function changeUserType(req: components.ChangeUserTypeRequest) {
	return webapi.post<components.AdminRes>(`/v1/admin/user/changetype`, req)
}

/**
 * @description "解禁用户"
 * @param req
 */
export function unBanUser(req: components.ChangeUserTypeRequest) {
	return webapi.post<components.AdminRes>(`/v1/admin/user/unban`, req)
}

/**
 * @description "获取评论列表"
 * @param params
 */
export function getCommentList(params: components.CommentListRequestParams) {
	return webapi.get<components.CommentListResponse>(`/v1/comment/list`, params)
}

/**
 * @description "点赞评论"
 * @param req
 */
export function likeComment(req: components.LikeCommentRequest) {
	return webapi.post<components.LikeCommentResponse>(`/v1/comment/like`, req)
}

/**
 * @description "创建评论"
 * @param params
 * @param req
 * @param headers
 */
export function newComment(params: components.NewCommentRequestParams, req: components.NewCommentRequest, headers: components.NewCommentRequestHeaders) {
	return webapi.post<components.NewCommentResponse>(`/v1/comment/new`, params, req, headers)
}

/**
 * @description "帧（文章）详情"
 * @param params
 */
export function getFeedDetail(params: components.GetFeedDetailRequestParams, id: number) {
	return webapi.get<components.GetFeedDetailResponse>(`/v1/feed/${id}`, params)
}

/**
 * @description "收藏帧（文章）"
 * @param req
 */
export function collectFeed(req: components.DeleteFeedRequest) {
	return webapi.post<components.StatusResponse>(`/v1/feed/collect`, req)
}

/**
 * @description "删除帧（文章）"
 * @param req
 */
export function deleteFeed(req: components.DeleteFeedRequest) {
	return webapi.post<components.StatusResponse>(`/v1/feed/delete`, req)
}

/**
 * @description "创作中心获取帧（文章）"
 * @param params
 */
export function getFeedListSelf(params: components.GetFeedListRequestParams) {
	return webapi.get<components.GetFeedListResponse>(`/v1/feed/feedlist/get`, params)
}

/**
 * @description "点赞帧（文章）"
 * @param req
 */
export function likeFeed(req: components.DeleteFeedRequest) {
	return webapi.post<components.StatusResponse>(`/v1/feed/like`, req)
}

/**
 * @description "创建帧（文章）"
 * @param params
 * @param req
 * @param headers
 */
export function newFeed(params: components.NewFeedRequestParams, req: components.NewFeedRequest, headers: components.NewFeedRequestHeaders) {
	return webapi.post<components.StatusResponse>(`/v1/feed/new`, params, req, headers)
}

/**
 * @description "分享小红书帧（文章）"
 */
export function shareFeedXHS() {
	return webapi.get<components.ShareFeedXHSResponse>(`/v1/feed/share/xhs`)
}

/**
 * @description "获取用户的粉丝列表"
 * @param req
 */
export function getFansList(req: components.FollowRequest) {
	return webapi.get<components.FollowListResponse>(`/v1/follow/fansList`, req)
}

/**
 * @description "关注/取消关注用户"
 * @param req
 */
export function follow(req: components.FollowRequest) {
	return webapi.post<components.FollowResponse>(`/v1/follow/follow`, req)
}

/**
 * @description "获取用户的关注列表"
 * @param req
 */
export function getFollowList(req: components.FollowRequest) {
	return webapi.get<components.FollowListResponse>(`/v1/follow/followList`, req)
}

/**
 * @description "获取首页信息流"
 * @param params
 */
export function getIndexFeeds(params: components.GetIndexFeedsRequestParams) {
	return webapi.get<components.GetIndexFeedsResponse>(`/v1/home/getfeeds`, params)
}

/**
 * @description "获取用户主页帖子"
 * @param params
 */
export function userFeeds(params: components.UserFeedsRequestParams) {
	return webapi.get<components.GetIndexFeedsResponse>(`/v1/home/userfeeds`, params)
}

/**
 * @description "举报"
 * @param req
 */
export function report(req: components.ReportRequest) {
	return webapi.post<components.ReportResponse>(`/v1/report/new-report`, req)
}

/**
 * @description "创建帧（文章）"
 * @param req
 */
export function normalTest(req: components.TestRequest) {
	return webapi.post<components.TestResponse>(`/v1/test/normal`, req)
}

/**
 * @description "获取到上传token"
 * @param params
 */
export function getUpToken(params: components.GetUpTokenRequestParams) {
	return webapi.get<components.GetUpTokenResponse>(`/v1/upload/normal`, params)
}

/**
 * @description "获取用户信息"
 * @param params
 */
export function userInfo(params: components.UserInfoRequestParams, uid: number) {
	return webapi.get<components.UserInfoResponse>(`/v1/user/${uid}`, params)
}

/**
 * @description "uid账号登录"
 * @param req
 */
export function login(req: components.LoginRequest) {
	return webapi.post<components.LoginResponse>(`/v1/user/login`, req)
}

/**
 * @description "手机密码登录"
 * @param params
 * @param req
 * @param headers
 */
export function loginByMobilePass(params: components.LoginMobilePassRequestParams, req: components.LoginMobilePassRequest, headers: components.LoginMobilePassRequestHeaders) {
	return webapi.post<components.LoginResponse>(`/v1/user/login-mobile-pass`, params, req, headers)
}

/**
 * @description "用户注册"
 * @param req
 */
export function register(req: components.RegisterRequest) {
	return webapi.post<components.RegisterResponse>(`/v1/user/register`, req)
}

/**
 * @description "获取验证码"
 * @param req
 */
export function verifyCode(req: components.VerifyCodeRequest) {
	return webapi.post<components.VerifyCodeResponse>(`/v1/user/verify-code`, req)
}

/**
 * @description "更改用户信息(token)"
 * @param req
 */
export function changeUserInfo(req: components.ChangeUserInfoRequest) {
	return webapi.post<components.ChangeUserInfoResponse>(`/v1/user/change-info`, req)
}

/**
 * @description "更改手机号码(token)"
 * @param req
 */
export function changeMobile(req: components.ChangeMobileRequest) {
	return webapi.post<components.ChangeMobileResponse>(`/v1/user/change-mobile`, req)
}

/**
 * @description "更改密码(token)"
 * @param req
 */
export function changePassword(req: components.ChangePasswordRequest) {
	return webapi.post<components.ChangePasswordResponse>(`/v1/user/change-password`, req)
}

/**
 * @description "获取用户关注、收藏、点赞列表"
 */
export function userRelation() {
	return webapi.get<components.UserRelationResponse>(`/v1/user/relation`)
}
