package gosdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SignSDK struct {
	AppID     string
	SecretKey string
	Timeout   int64 // 时间戳超时时间，单位：秒
}

// 必备签名参数
type NeededSignParams map[string]string

func (n NeededSignParams) ToJson() string {
	signParamsBody, _ := json.Marshal(n)
	return string(signParamsBody)
}

func (n NeededSignParams) ToString() string {
	var keyArr []string
	for k := range n {
		if v := n[k]; len(v) == 0 {
			continue
		}
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	var keyValueArray []string
	for _, key := range keyArr {
		keyValueArray = append(keyValueArray, key+"="+n[key])
	}
	return strings.Join(keyValueArray, "&")
}

const APPInfoInKey = "_jx3box_ak_"
const SignResultKey = "_jx3box_sign_"

// 时间戳校验
// 正负5分钟均可通过
func (s *SignSDK) isLegalTime(timeRaw string) bool {
	timeoutMax := int64(5 * 60)
	if s.Timeout > 0 {
		timeoutMax = s.Timeout
	}

	i, err := strconv.ParseInt(timeRaw, 10, 64)
	if err != nil {
		return false
	}
	return math.Abs(float64(time.Now().Unix()-i)) < float64(timeoutMax)
}

// 签名算法
// @params sk <string> 密钥
// @params urlParams <url.Values> 待签名计算参数, 其中 _jx3box_ak_ 参数中包含了 appid 和 nonce_str
// @return string 签名结果
func (s *SignSDK) signParams(sk string, urlParams url.Values, appInfo NeededSignParams) string {
	// 计算url原有参数
	var keyArr []string
	for k := range urlParams {
		if k == APPInfoInKey || k == SignResultKey {
			continue
		}
		if v := urlParams[k]; len(v) == 0 {
			continue
		}
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	var keyValueArray []string
	for _, key := range keyArr {
		keyValueArray = append(keyValueArray, key+"="+strings.Join(urlParams[key], ","))
	}
	// 合并url原有参数，app信息，密钥，然后计算sign值
	keyValueArray = append(keyValueArray, appInfo.ToString(), "sk="+sk)
	beSignStr := strings.Join(keyValueArray, "&")
	hasher := md5.New()
	hasher.Write([]byte(beSignStr))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

// 获取签名后的url
func (s *SignSDK) GetSignedURL(api string) (string, error) {
	urlObj, err := url.Parse(api)
	if err != nil {
		return "", err
	}

	query := urlObj.Query()

	signParams := NeededSignParams{
		"appid":   s.AppID,
		"non_str": randomNonceStr(10),
		"t":       strconv.FormatInt(time.Now().Unix(), 10),
	}

	query.Set(APPInfoInKey, signParams.ToJson())
	sign := s.signParams(s.SecretKey, query, signParams)
	query.Set(SignResultKey, sign)
	query.Set(APPInfoInKey, signParams.ToString())
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (s *SignSDK) CheckSign(beCheckSign string, urlParams url.Values) bool {
	appInfo := urlParams.Get(APPInfoInKey)
	if appInfo == "" {
		return false
	}
	var signParams NeededSignParams
	if err := json.Unmarshal([]byte(appInfo), &signParams); err != nil {
		return false
	}
	// 时间戳校验
	if !s.isLegalTime(signParams["t"]) {
		return false
	}
	// 签名校验
	sign := s.signParams(s.SecretKey, urlParams, signParams)
	return sign != "" && beCheckSign == sign
}
