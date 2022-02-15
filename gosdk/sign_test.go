package gosdk

import (
	"log"
	"net/url"
	"testing"
)

func TestSDK(t *testing.T) {

	api := "http://node.jx3box.com/skills"
	sdk := SignSDK{
		AppID:     "test-jx3box",
		SecretKey: "test-jx3box-2021-hello",
	}

	resourceURL2, _ := sdk.GetSignedURL(api)
	log.Println(resourceURL2)

	urlObj, _ := url.Parse(resourceURL2)
	query := urlObj.Query()

	log.Println(sdk.CheckSign(query.Get("_jx3box_sign_"), query))

	// var client = http.Client{}
	// request2, _ := http.NewRequest("GET", resourceURL2, nil)

	// response2, err := client.Do(request2)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer response2.Body.Close()
	// body2, _ := ioutil.ReadAll(response2.Body)
	// if response2.StatusCode == 200 {
	// 	t.Log(string(body2))
	// } else {
	// 	t.Logf("code: %d body: %s", response2.StatusCode, string(body2))
	// }
}
