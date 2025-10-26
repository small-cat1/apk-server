package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	projectModel "ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	projectRes "ApkAdmin/model/project/response"
	systemRes "ApkAdmin/model/system/response"
	"ApkAdmin/service/project"
	"ApkAdmin/utils"
	"ApkAdmin/utils/captcha"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"time"
)

var store = captcha.NewDefaultRedisStore()

type BaseApi struct {
}

func (a *BaseApi) CustomerServiceConfig(c *gin.Context) {
	config, err := websiteConfigService.GetConfig("service")
	if err != nil {
		global.GVA_LOG.Error("客服配置获取失败!", zap.Error(err))
		response.FailWithMessage("客服配置获取失败", c)
		return
	}
	helper := utils.MapHelper(config)
	resp := projectRes.CustomerServiceConfigResp{
		Enabled:    helper.GetBool("enabled"),
		ShowText:   true,
		ButtonText: helper.GetString("button_text", "客服"),
		Tooltip:    helper.GetString("tooltip", "联系客服"),
		Position: projectRes.CustomerServicePosition{
			Right:  "20px",
			Bottom: "80px",
			ZIndex: 999,
		},
		Contacts: projectRes.CustomerServiceContacts{
			Qq:           helper.GetString("qq"),
			Wechat:       helper.GetString("wechat"),
			WechatQrcode: helper.GetString("wechat_qrcode"),
			Phone:        helper.GetString("phone"),
			Email:        helper.GetString("email"),
			Im:           helper.GetBool("im_switch"),
		},
		WorkTime: helper.Get("work_time"),
		Notice:   helper.Get("notice"),
		ImType:   "",
		ImConfig: projectRes.ImConfig{},
		Preload:  false,
	}
	response.OkWithData(resp, c)
}

// Captcha 获取验证码
func (a *BaseApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}
	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(
		96,  // 高度：32 * 2
		300, // 宽度：100 * 2
		4,   // 验证码长度
		0.7, // 最大倾斜度
		80,  // 干扰点数量
	)
	cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c)) // v8下使用redis
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.GVA_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GVA_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证码获取成功", c)
}

// Register 用户注册
func (a BaseApi) Register(c *gin.Context) {
	var req request.BaseRegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("注册失败，参数不正确!", zap.Error(err))
		response.FailWithMessage("注册失败，参数不正确", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage("注册失败，错误："+err.Error(), c)
		return
	}
	clientIP := c.ClientIP()
	// ==================== 图形验证码验证 ====================
	// 所有注册都必须验证图形验证码（不像登录可以前几次不验证）
	if req.Captcha == "" || req.CaptchaKey == "" {
		incrementIPAttempt("register", clientIP)
		response.FailWithMessage("请输入图形验证码", c)
		return
	}
	if !store.Verify(req.CaptchaKey, req.Captcha, true) {
		incrementIPAttempt("register", clientIP)
		response.FailWithMessage("图形验证码错误", c)
		return
	}
	// 检查手机号注册频率（防止用同一手机号反复注册）
	//if err := checkPhoneRegisterFrequency(req.Phone); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	// 检查手机号是否已注册
	if userExists, _ := checkPhoneExists(req.Phone); userExists {
		response.FailWithMessage("该手机号已被注册", c)
		return
	}
	// 邮箱验证
	if emailExists, _ := checkEmailExists(req.Email); emailExists {
		response.FailWithMessage("该邮箱已被注册", c)
		return
	}
	err = UserService.RegisterUser(req, clientIP)
	if err != nil {
		incrementIPAttempt("register", clientIP)
		global.GVA_LOG.Error("注册失败，", zap.Error(err))
		response.FailWithMessage("注册失败，请稍后再试！", c)
		return
	}
	// 注册成功，清除限制记录
	clearIPAttempt("register", clientIP)
	response.OkWithMessage("注册成功", c)

}

// Login 用户登录
func (a BaseApi) Login(c *gin.Context) {
	var req request.BaseLoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("登录失败，参数不正确!", zap.Error(err))
		response.FailWithMessage("登录失败，参数不正确", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage("登录失败，错误："+err.Error(), c)
		return
	}
	clientIP := c.ClientIP()
	// ==================== IP 级别限制检查 ====================
	if isBlocked, until := checkIPBlocked(clientIP); isBlocked {
		response.FailWithMessage(fmt.Sprintf("该IP已被临时封禁至 %s", until.Format("15:04")), c)
		return
	}
	// ==================== 图形验证码验证 ====================
	// 所有注册都必须验证图形验证码（不像登录可以前几次不验证）
	if req.Captcha == "" || req.CaptchaKey == "" {
		incrementIPAttempt("login", clientIP)
		response.FailWithMessage("请输入图形验证码", c)
		return
	}
	if !store.Verify(req.CaptchaKey, req.Captcha, true) {
		incrementIPAttempt("login", clientIP)
		response.FailWithMessage("图形验证码错误", c)
		return
	}
	u := &projectModel.User{
		Phone:        &req.Phone,
		PasswordHash: req.Password,
	}
	userAgent := c.Request.UserAgent()
	user, err := UserService.Login(u, clientIP, userAgent)
	if err != nil {
		global.GVA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		// 验证码次数+1
		incrementIPAttempt("login", clientIP)
		response.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if !user.AccountStatus.IsActive() {
		global.GVA_LOG.Error("登陆失败! 用户被禁止登录!")
		// 验证码次数+1
		incrementIPAttempt("login", clientIP)
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	//登录以后签发jwt
	a.TokenNext(c, *user, clientIP)
}

// TokenNext 登录以后签发jwt
func (a *BaseApi) TokenNext(c *gin.Context, user projectModel.User, clientIP string) {
	token, claims, err := utils.ClientLoginToken(&user)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	clearIPAttempt("login", clientIP)
	response.OkWithDetailed(projectRes.LoginResponse{
		User:      user,
		Token:     token,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	}, "登录成功", c)
	return
}

// ==================== 防护函数 ====================

// 检查手机号注册频率
func checkPhoneRegisterFrequency(phone string) error {
	key := "register:phone:" + phone
	// 防止用同一手机号反复尝试注册（可能是在测试或攻击）
	v, ok := global.BlackCache.Get(key)
	if ok {
		lastAttempt := v.(time.Time)
		// 同一手机号5分钟内只能尝试注册一次
		if time.Since(lastAttempt) < 5*time.Minute {
			remain := int(5*time.Minute - time.Since(lastAttempt))
			return fmt.Errorf("请%d秒后再试", remain)
		}
	}
	// 记录本次尝试时间
	global.BlackCache.Set(key, time.Now(), 5*time.Minute)
	return nil
}

// checkIPBlocked 检查 IP 是否被封禁（新增辅助函数）
func checkIPBlocked(ip string) (bool, time.Time) {
	blockKey := "register:block:" + ip
	if v, ok := global.BlackCache.Get(blockKey); ok {
		if blockUntil, ok := v.(time.Time); ok {
			if time.Now().Before(blockUntil) {
				return true, blockUntil
			}
		}
	}
	return false, time.Time{}
}

// 累加IP失败次数
func incrementIPAttempt(action, ip string) {
	key := action + ":attempt:" + ip
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Hour)
	} else {
		newCount := interfaceToInt(v) + 1
		global.BlackCache.Set(key, newCount, time.Hour)
		// 失败次数过多，临时封禁
		if newCount >= 20 {
			blockKey := action + ":block:" + ip
			blockUntil := time.Now().Add(time.Hour)
			global.BlackCache.Set(blockKey, blockUntil, time.Hour)
		}
	}
}

// 清除IP尝试记录
func clearIPAttempt(action, ip string) {
	key := action + ":attempt:" + ip
	global.BlackCache.Delete(key)
}

// 检查手机号是否已存在
func checkPhoneExists(phone string) (bool, error) {
	// 查询数据库
	count, err := UserService.CountUserBy(project.WithPhone(phone))
	return count > 0, err
}

// 检查邮箱是否已存在
func checkEmailExists(email string) (bool, error) {
	count, err := UserService.CountUserBy(project.WithEmail(email))
	return count > 0, err
}

// 类型转换
func interfaceToInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch v.(type) {
	case int:
		return v.(int)
	case int64:
		return int(v.(int64))
	case float64:
		return int(v.(float64))
	default:
		return 0
	}
}
