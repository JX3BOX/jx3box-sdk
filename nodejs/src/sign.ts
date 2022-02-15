import * as crypto from "crypto"
import { URLSearchParams, URL } from "url"
import { random as randomStr } from "./utils"
const APPInfoInKey = "_jx3box_ak_"
const SignResultKey = "_jx3box_sign_"

class SignSDK {
    private appid: string;
    private secretKey: string;
    private timeout: number;// 时间误差
    constructor(appid: string, secretKey: string, timeout: number = 5 * 60) {
        this.appid = appid;
        this.secretKey = secretKey;
        this.timeout = timeout
    }
    public checkSign(beCheckSign: string | null, urlParams: URLSearchParams): boolean {
        const appInfo = urlParams.get(APPInfoInKey)
        if (!appInfo) {
            return false
        }
        let appInfoObj: { [key: string]: string } = {}
        try {
            appInfoObj = JSON.parse(appInfo)
        } catch (e) {
            return false
        }
        if (!this.isLegalTime(appInfoObj["t"])) {
            return false
        }
        const sign = this.sign(urlParams, appInfoObj)
        return sign !== "" && sign === beCheckSign
    }
    private isLegalTime(timeRaw: string): boolean {
        const time = parseInt(timeRaw)
        if (isNaN(time)) {
            return false
        }
        const now = Date.now() / 1000
        return Math.abs(now - time) < this.timeout
    }
    private sign(params: URLSearchParams, appInfo: { [key: string]: string }): string {
        const appInfoArr = Object.keys(appInfo).sort()
        const appInfoStr = appInfoArr.map(key => `${key}=${appInfo[key]}`).join("&")

        const keyMaps: { [key: string]: boolean } = {}
        for (let key of params.keys()) {
            if (key == APPInfoInKey || key == SignResultKey) {
                continue
            }
            keyMaps[key] = true
        }
        const keys = Object.keys(keyMaps).sort()
        const values: Array<string> = keys.map((key: string) => { return key + '=' + params.getAll(key).join(",") })
        values.push(appInfoStr)
        values.push("sk=" + this.secretKey)
        const beSignStr = values.join("&")
        const sign = crypto.createHash('md5').update(beSignStr).digest("hex");
        return sign.toUpperCase()

    }

    getSignedURL(urlRaw: string): string {
        const urlObj = new URL(urlRaw)

        const search = urlObj.searchParams

        const appInfo = {
            "appid": this.appid,
            "non_str": randomStr(10),
            "t": ((Date.now() / 1000).toFixed()).toString()
        }
        const sign = this.sign(search, appInfo)
        search.set(APPInfoInKey, JSON.stringify(appInfo))
        search.set(SignResultKey, sign)
        urlObj.search = search.toString()
        return urlObj.toString()
    }
}

export default SignSDK