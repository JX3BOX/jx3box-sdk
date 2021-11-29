import { Sign as SignSDK } from "../index"
import * as http from "http"


function TestSDK() {
    const targetURL = "http://gatway.test.com/xxxx?a=1&xxx=2&xxx=23123"
    const sdk = new SignSDK("a", "b")

    const raw = sdk.getSignedURL(targetURL)
    console.log(raw)

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