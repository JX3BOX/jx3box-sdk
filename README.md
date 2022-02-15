## JX3BOX SDK

- 获取jx3box资源
- SSO

### 语言支持

- golang
- nodejs

### Demo

见 `gosdk/sign_test.go` 或 `nodejs/src/test/*.test`


## 签名算法
该签名算法已经实现了golang,nodejs版本。其他语言实现请参 签名规则说明

### 签名规则说明

#### 通用参数说明

```
appid {string} 应用id
nonce_str {string} 随机字符串
t {int64} 时间戳 单位 秒
sign {string} 签名
```

#### 步骤1 配置通用参数

将 `appid`, `nonce_str`, `t`拼接成json得到结果:【注意时间戳也需要转成字符串】

```json
{"appid":"abc", "nonce_str":"random_xxxx", "t":"179999999"}
```

记为`_jx3box_ak_`,

然后将  `appid`, `nonce_str`, `t` 进行字典排序，如字典排序后的结果为

```json
["appid", "nonce_str", "t"]
```

依次将`key`,`value` 用`=` 连接，然后用`&`拼接结果：

```plain
appid=abc&nonce_str=random_xxxx&t=179999999
```

记为 `appBeSignIno`


#### 步骤2 从url参数中获取待签名的字符串。

假如目标接口为：【标记1】

```
https://open.jx3box.com.com/item-price?username=c&item_name=1x
```

将所有url查询参数进行字典排序得到：

```
["item_name", "username"]
```

依次将`key`,`value` 用`=` 连接，然后用`&`拼接结果：

```plain
item_name=1x&username=c
```

然后再将拼接结果后面加上 `appBeSignIno`的值和`sk=私钥`, 同样使用`&`拼接 得到【假设私钥是123456】：

```javascript
let beSiginStr = ["item_name=1x&username=c","appid=abc&nonce_str=random_xxxx",  "sk=123456"].join("&")
// ===>appid=abc&nonce_str=random_xxxx&t=179999999&sk=123456&item_name=1x&username=c
// 如果 目标地址的查询字符串为空，那么待签名字段为:  ["appid=abc&nonce_str=random_xxxx",  "sk=123456"].join("&")
```

这样就得到需要签名的字符串。将该字符串进行 `md5hex` 计算，最后得到 签名结果: `XXXXXXXXXXXXXXXXX`

将 签名结果作为参数`_jx3box_sign_`的值 和 步骤一中的`_jx3box_ak_`拼接到实际接口后，即完成了接口签名


因此最终的访问url为【可以对比标记1】:

```
https://open.jx3box.com.com/item-price?username=c&item_name=1x&_jx3box_ak_={"appid":"abc", "nonce_str":"random_xxxx", "t":"179999999"}&_jx3box_sign_=XXXXXXXXXXXXXXXXX
```