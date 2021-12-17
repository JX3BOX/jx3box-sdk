import * as crypto from "crypto"
import { URLSearchParams, URL } from "url"

class SignSDK {
    private appid: string;
    private secretKey: string;
    constructor(appid: string, secretKey: string) {
        this.appid = appid;
        this.secretKey = secretKey;
    }
    public checkSign(beCheckSign: string | null, urlParams: URLSearchParams): boolean {
        const sign = this.signParams(urlParams)
        return sign !== "" && sign === beCheckSign
    }
    private signParams(params: URLSearchParams): string {
        const appInfo = params.get("__ak__")
        if (!appInfo) {
            return ""
        }
        const appInfoObj = JSON.parse(appInfo)
        const appInfoArr = Object.keys(appInfoObj).sort()
        const appInfoStr = appInfoArr.map(key => `${key}=${appInfoObj[key]}`).join("&")

        const keyMaps: { [key: string]: boolean } = {}
        for (let key of params.keys()) {
            if (key == "__ak__" || key == "sign") {
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
    private random(length: number): string {
        var result = '';
        var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        var charactersLength = characters.length;
        for (var i = 0; i < length; i++) {
            result += characters.charAt(Math.floor(Math.random() *
                charactersLength));
        }
        return result
    }
    getSignedURL(urlRaw: string): string {
        const urlObj = new URL(urlRaw)

        const search = urlObj.searchParams

        const signParams = {
            "appid": this.appid,
            "non_str": this.random(10),
            "t": (Date.now() / 1000).toFixed()
        }

        search.set("__ak__", JSON.stringify(signParams))

        const sign = this.signParams(search)

        search.set("sign", sign)

        urlObj.search = search.toString()

        return urlObj.toString()
    }
}

export default SignSDK