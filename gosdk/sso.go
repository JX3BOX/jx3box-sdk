package gosdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SSO struct {
	AppID     string `yaml:"appid"`
	SecretKey string `yaml:"secretKey"`
	API       string `yaml:"api"`
}

func (s *SSO) GetLoginPage(params map[string]string, backTo ...string) (string, error) {

	urlObj, err := url.Parse(s.API)
	if err != nil {
		return "", err
	}
	query := urlObj.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	if len(backTo) > 0 {
		query.Set("__backto__", backTo[0])
	}
	urlObj.Path = "/authorize"
	urlObj.RawQuery = query.Encode()

	sign := SignSDK{
		AppID:     s.AppID,
		SecretKey: s.SecretKey,
	}
	return sign.GetSignedURL(urlObj.String())

}

func (s *SSO) GetResource(scope string, token string, data interface{}) error {

	urlObj, err := url.Parse(s.API)
	if err != nil {
		return err
	}
	urlObj.Path = "/authorize/resource"
	query := urlObj.Query()
	query.Set("scope", scope)
	query.Set("resource_token", token)
	urlObj.RawQuery = query.Encode()
	sign := SignSDK{
		AppID:     s.AppID,
		SecretKey: s.SecretKey,
	}

	targetURL, err := sign.GetSignedURL(urlObj.String())
	if err != nil {
		return err
	}

	response, err := http.Get(targetURL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("status code: %d message: %s", response.StatusCode, string(body))
	}

	return json.Unmarshal(body, data)
}
