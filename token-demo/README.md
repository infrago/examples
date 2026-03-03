# token-demo

启动：

```bash
go run ./token-demo
```

签发 token：

```bash
curl -s -X POST 'http://127.0.0.1:8096/token/sign?uid=u2001'
```

使用 token 访问签名路由：

```bash
TOKEN='上一步返回的token'
curl -s -H "Authorization: Bearer $TOKEN" 'http://127.0.0.1:8096/token/profile'
```

访问 auth 路由：

```bash
TOKEN='上一步返回的token'
curl -s -H "Authorization: Bearer $TOKEN" 'http://127.0.0.1:8096/token/auth'
```

吊销 tokenId：

```bash
TOKEN='上一步返回的token'
curl -s -X POST -H "Authorization: Bearer $TOKEN" 'http://127.0.0.1:8096/token/revoke'
```

再次访问 profile/auth 将失败（未签名/未认证）。
