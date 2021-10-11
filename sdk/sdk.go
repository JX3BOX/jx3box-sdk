package sdk

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

type SignSDK struct {
	AppID     string
	SecretKey string
}

func (s *SignSDK) signParams(sk string, urlParams url.Values) string {
	var keyArr []string

	for k := range urlParams {
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
	keyValueArray = append(keyValueArray, "sk="+sk)
	beSignStr := strings.Join(keyValueArray, "&")
	log.Println(beSignStr)
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

	query.Add("appid", s.AppID)
	query.Add("nonce_str", s.randomNonceStr(10))
	query.Add("__t", fmt.Sprintf("%d", time.Now().Unix()))
	sign := s.signParams(s.SecretKey, query)
	query.Add("sign", sign)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (s *SignSDK) CheckSign(beCheckSign string, urlParams url.Values) bool {
	return beCheckSign == s.signParams(s.SecretKey, urlParams)
}
