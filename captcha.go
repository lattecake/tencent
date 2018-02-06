package tencent

func New(appId string, secretId string, secretKey string, proxy string) CaptchaService {

	Init(appId, secretId, secretKey, proxy)

	return NewService(getServiceMiddleware(), config)
}

func getServiceMiddleware() (mw []Middleware) {
	mw = []Middleware{}
	// Append your middleware here

	return
}
