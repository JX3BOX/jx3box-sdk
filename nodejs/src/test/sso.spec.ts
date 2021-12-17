import { equal } from "assert";
import { URL } from "url"
import { SSO } from "../index"

describe("SignSDK", () => {
    it("获取登录地址", () => {
        const sso = new SSO({ appid: "a", secretKey: "b" })
        console.log(sso.getLoginPage({ a: "1", b: "2" }, "http://localhost:3000"))
    });
});
