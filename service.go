package tencent

import (
	"math/rand"
	"strconv"
	"time"
	"encoding/base64"
	"net/url"
	"strings"
	"sort"
	"encoding/json"
	"github.com/lattecake/request"
	"errors"
	"crypto/sha1"
	"crypto/hmac"
)

type CaptchaService interface {
	CaptchaIframeQuery(ip string, captchaType string, level string, clientType string, businessId string) (url string, err error)
	CaptchaCheck(ticket string, ip string, captchaType string, verifyType string, businessId string) (res result, err error)
}

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Url     string `json:"url"`
	IsRight int    `json:"is_right"`
}

const URL = "csec.api.qcloud.com/v2/index.php"

type basicCaptchaService struct {
	config Config
}

func newBasicCaptchaService(config Config) (s *basicCaptchaService) {
	return &basicCaptchaService{
		config: config,
	}
}

func NewService(middleware []Middleware, config Config) CaptchaService {
	var svc CaptchaService = newBasicCaptchaService(config)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func (c *basicCaptchaService) CaptchaIframeQuery(ip string, captchaType string, level string, clientType string, businessId string) (url string, err error) {

	if captchaType == "" {
		captchaType = GetConf().CaptchaType
	}
	if level == "" {
		level = GetConf().DisturbLevel
	}
	if clientType == "" {
		clientType = GetConf().ClientType
	}

	if businessId == "" {
		businessId = GetConf().BusinessId
	}

	httpUrl := makeURL("GET", "CaptchaIframeQuery", GetConf().Region, GetConf().SecretId, GetConf().SecretKey, map[string]string{
		"userIp":       ip,
		"captchaType":  captchaType,
		"disturbLevel": level,
		"isHttps":      "1",
		"clientType":   clientType,
		"businessId":   businessId,
		"appId":        GetConf().AppId,
	})

	body, err, _ := request.Get(httpUrl, nil, nil, GetConf().Proxy)
	if err != nil {
		return
	}

	var res result

	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	if res.Code != 0 {
		return "", errors.New(res.Message)
	}

	return res.Url, nil
}

func (c *basicCaptchaService) CaptchaCheck(ticket string, ip string, captchaType string, verifyType string, businessId string) (res result, err error) {

	if captchaType == "" {
		captchaType = GetConf().CaptchaType
	}

	if verifyType == "" {
		verifyType = GetConf().VerifyType
	}

	if businessId == "" {
		businessId = GetConf().BusinessId
	}

	httpUrl := makeURL("GET", "CaptchaCheck", GetConf().Region, GetConf().SecretId, GetConf().SecretKey, map[string]string{
		"userIp":      ip,
		"captchaType": captchaType,
		"ticket":      ticket,
		"verifyType":  verifyType,
		"businessId":  businessId,
	})

	body, err, _ := request.Get(httpUrl, nil, nil, GetConf().Proxy)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}

	if res.Code != 0 {
		return res, errors.New(res.Message)
	}

	return
}

func makeURL(method string, action string, region string, secretId string, secretKey string, args map[string]string) string {

	args["Nonce"] = strconv.Itoa(rand.Intn(0x7fffffff))
	args["Action"] = action
	args["Region"] = region
	args["SecretId"] = secretId
	args["Timestamp"] = strconv.Itoa(int(time.Now().Unix()))

	args["Signature"] = base64.StdEncoding.EncodeToString([]byte(common.HashHMac(method+URL+"?"+makeQueryString(args, false), secretKey)))

	return "https://" + URL + "?" + makeQueryString(args, true)
}

func makeQueryString(args map[string]string, isURLEncoded bool) string {

	var res []string

	var keys []string
	for k := range args {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		if !isURLEncoded {
			res = append(res, v+"="+args[v])
		} else {
			res = append(res, v+"="+url.QueryEscape(args[v]))
		}
	}

	return strings.Join(res, "&")
}

func hashHMac(input string, secret string) string {
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	_, err := mac.Write([]byte(input))
	if err != nil {
		return err.Error()
	}
	return string(mac.Sum(nil))
}
