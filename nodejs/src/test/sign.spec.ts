import { equal } from "assert";
import { URL } from "url"
import { SignSDK } from "../index"

describe("SignSDK", () => {
    it("签名与校验", () => {
        const targetURL = "http://gatway.test.com/xxxx?a=1&xxx=2&xxx=23123"
        const sdk = new SignSDK("a", "b")
        const raw = sdk.getSignedURL(targetURL)
        console.log("签名前的目标地址为:", targetURL)
        console.log("签名后的目标地址为:", raw)

        const urlObj = new URL(raw)
        const sdkError = new SignSDK("a", "c")
        equal(sdk.checkSign(urlObj.searchParams.get("_jx3box_sign_"), urlObj.searchParams), true);
        equal(sdkError.checkSign(urlObj.searchParams.get("_jx3box_sign_"), urlObj.searchParams), false);
    });
});
