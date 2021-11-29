import * as crypto from "crypto"
import * as http from "http"
import { URLSearchParams, URL } from "url"

class SignSDK {
    private appid: string;
    private secretKey: string;
    constructor(appid: string, secretKey: string) {
        this.appid = appid;
        this.secretKey = secretKey;
    }
    private signParams(params: URLSearchParams): string {
        const keyMaps: { [key: string]: boolean } = {}
        for (let key of params.keys()) {
            keyMaps[key] = true
        }
        const keys = Object.keys(keyMaps)
        keys.sort()
        const values: Array<string> = []
        keys.forEach((key: string) => { values.push(key + '=' + params.getAll(key).join(",")) })
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
        search.set("appid", this.appid)
        search.set("nonce_str", this.random(10))
        search.set("__t", (Date.now() / 1000).toFixed())

        const sign = this.signParams(search)

        search.set("sign", sign)

        urlObj.search = search.toString()

        return urlObj.toString()
    }
}

export default SignSDK