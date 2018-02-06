package tencent

type Config struct {
	Proxy        string
	SecretId     string
	SecretKey    string
	CaptchaType  string // 验证码类型
	DisturbLevel string // 验证码难度
	ClientType   string // 客户端类型
	BusinessId   string
	VerifyType   string
	Region       string
	AppId        string
}

var config Config

func Init(appId string, secretId string, secretKey string, proxy string /*, captchaType string, level string, clientType string, businessId string, verifyType string*/) {

	config.SecretId = secretId
	config.SecretKey = secretKey
	config.CaptchaType = "4"
	config.DisturbLevel = "1"
	config.ClientType = "2"
	config.BusinessId = "1256065690"
	config.VerifyType = "2"
	config.Region = "bj"
	config.AppId = appId

	if proxy != "" {
		config.Proxy = proxy
	}
}

func GetConf() *Config {
	return &config
}
