# tencent
腾讯验证码服务

`go get github.com/lattecake/tencent`

## 使用

```go
QQCaptcha = tencent.New("your tencentAppId", "your secretId", "your secretKey", "your http_proxy")

```

### 生成

```go
qqUrl, err := captcha.QQCaptcha.CaptchaIframeQuery("client ip addr", "9", "1", "1", "your tencent appId")
```


### 验证

```go

res, err := captcha.QQCaptcha.CaptchaCheck("client ticket", "client ip addr", "9", "2", "your tencentAppId")

if err != nil {
    fmt.Println(err)
    return
}

if res.Code != 0 {
    fmt.Prineln(res.Message)
    return
}

if res.IsRight != 1 {
    return
}
```
