package sdk

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestSDK(t *testing.T) {

	api := "http://gatway.test.com/xxxx?a=1&xxx=2&xxx=23123"
	sdk := SignSDK{
		AppID:     "a",
		SecretKey: "b",
	}

	resourceURL2, _ := sdk.GetSignedURL(api)
	log.Println(resourceURL2)
	var client = http.Client{}
	request2, _ := http.NewRequest("GET", resourceURL2, nil)

	response2, err := client.Do(request2)
	if err != nil {
		t.Fatal(err)
	}
	defer response2.Body.Close()
	body2, _ := ioutil.ReadAll(response2.Body)
	if response2.StatusCode == 200 {
		t.Log(string(body2))
	} else {
		t.Logf("code: %d body: %s", response2.StatusCode, string(body2))
	}
}
