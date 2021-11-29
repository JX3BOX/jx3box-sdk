## jx3box 资源访问 SDK

### API

```golang
package main

import (
    "github.com/JX3BOX/jx3box-sdk-golang/sdk"
    "net/http"
)


func main(){
    var signSDK =  sdk.SignSDK{
        AppID:     "a",
        SecretKey: "b",
    }
    resourceURL2, _ := signSDK.GetSignedURL("https://node.jx3box.com/xxxxx?Xxx=xxx")
    resp := http.Get(resourceURL2)
    ...
}


### 使用

```shell
go get -u github.com/JX3BOX/jx3box-sdk-golang/sdk
```

demo

```golang
package demo
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
```