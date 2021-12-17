package gosdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SignSDK struct {
	AppID     string
	SecretKey string
}

type NeededSignParams map[string]string

func (n NeededSignParams) ToJson() string {
	signParamsBody, _ := json.Marshal(n)
	return string(signParamsBody)
}

func (n NeededSignParams) ToString() string {
	var keyArr []string
	for k := range n {
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	var keyValueArray []string
	for _, key := range keyArr {
		keyValueArray = append(keyValueArray, key+"="+n[key])
	}
	return strings.Join(keyValueArray, "&")
}

// 签名算法
// ?a=x&b=y&c=z&__app={xxx:xxxxx}
// 对__app中的xxx:xxxxx 进行排序

func (s *SignSDK) signParams(sk string, urlParams url.Values) string {

	appInfo := urlParams.Get("__ak__")
	if appInfo == "" {
		return ""
	}
	var signParams NeededSignParams
	if err := json.Unmarshal([]byte(appInfo), &signParams); err != nil {
		return ""
	}

	var keyArr []string

	for k := range urlParams {
		if k == "__ak__" || k == "sign" {
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
	keyValueArray = append(keyValueArray, signParams.ToString(), "sk="+sk)
	beSignStr := strings.Join(keyValueArray, "&")
	hasher := md5.New()
	hasher.Write([]byte(beSignStr))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func (s *SignSDK) randomNonceStr(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func (s *SignSDK) GetSignedURL(api string) (string, error) {
	urlObj, err := url.Parse(api)
	if err != nil {
		return "", err
	}

	query := urlObj.Query()

	signParams := NeededSignParams{
		"appid":   s.AppID,
		"non_str": s.randomNonceStr(10),
		"t":       strconv.FormatInt(time.Now().Unix(), 10),
	}

	query.Set("__ak__", signParams.ToJson())
	sign := s.signParams(s.SecretKey, query)
	query.Set("sign", sign)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (s *SignSDK) CheckSign(beCheckSign string, urlParams url.Values) bool {
	sign := s.signParams(s.SecretKey, urlParams)
	return sign != "" && beCheckSign == sign
}
