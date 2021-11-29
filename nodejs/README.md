# jx3box SDK


## 安装

```shell
npm install @jx3box/sdk --save
```

## 快速开始

### 资源访问

```js
import {SignSDK} from "@jx3box/sdk"
import * as http from "http"


function TestSDK() {
    const appid = "A"
    const secretKey = "B"
    const targetURL = "http://resource.jx3box.com/xxxx?a=1&xxx=2&xxx=23123"
    const sdk = new SignSDK(appid, secretKey)
    const raw = sdk.getSignedURL(targetURL)
    http.get(raw, (res: http.IncomingMessage) => {
        res.setEncoding('utf8');
        let rawData = '';
        res.on('data', (chunk) => { rawData += chunk; });
        res.on('end', () => {
            try {
                console.log(rawData);
            } catch (e: any) {
                console.error(e.message);
            }
        });
    }).on('error', (e) => {
        console.error(`Got error: ${e.message}`);
    });
}

TestSDK()

```

### 第三方登录

