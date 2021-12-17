import { URL } from "url"
import Sign from "./sign"
import * as http from "http"
import * as https from "https"

export default class SSO {
    // app id
    private appid: string;
    // secret key 密钥
    private sk: string;
    // sso 服务器域名
    private api: string;
    constructor(settings: { appid: string; secretKey: string, api?: string }) {
        this.appid = settings.appid;
        this.sk = settings.secretKey;
        this.api = settings.api || "https://sso.jx3box.com";

    }
    // 获取sso登录页地址
    public getLoginPage(params: { [key: string]: string }, backto?: string): string {
        const urlObj = new URL(this.api)
        const query = urlObj.searchParams
        Object.keys(params).forEach(key => {
            query.set(key, params[key])
        })
        query.set("__backto__", backto || "")
        urlObj.pathname = "/authorize"
        urlObj.search = query.toString()
        const sign = new Sign(this.appid, this.sk)
        return sign.getSignedURL(urlObj.toString())
    }
    // 获取资源 如用户信息等
    public getResource(scope: string, token: string): Promise<any> {
        return new Promise<any>((resolve, reject) => {
            const urlObj = new URL(this.api)
            urlObj.pathname = "/authorize/resource"
            const query = urlObj.searchParams
            query.set("scope", scope)
            query.set("resource_token", token)
            urlObj.search = query.toString()
            const sign = new Sign(this.appid, this.sk)
            const targetURL = sign.getSignedURL(urlObj.toString())
            const responseHandle = (res: http.IncomingMessage) => {
                res.setEncoding('utf8');
                let rawData = '';
                res.on('data', (chunk: string) => { rawData += chunk; });
                res.on('end', () => {
                    if (res.statusCode != 200) {
                        reject(`status code: ${res.statusCode} error message: ${rawData}`);
                        return
                    }
                    try {
                        const parsedData = JSON.parse(rawData);
                        resolve(parsedData)
                    } catch (e) {
                        reject(`parse json error: ${rawData}`);
                    }
                });
            }
            const errorHandler = (err: Error) => {
                reject(`request error: ${err.message}`);
            }
            switch (urlObj.protocol) {
                case "https:":
                    https.get(targetURL, responseHandle).on('error', errorHandler);
                    break
                case "http:":
                    http.get(targetURL, responseHandle).on('error', errorHandler);
            }
        })
    }

}