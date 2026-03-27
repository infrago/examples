# ws-demo

最小可跑的 `web + ws` 示例，同时演示：

- `ws` 模块自动提供的默认 Upgrade 接管
- 自定义 `web.Endpoint`
- `ctx.Upgrade()`
- `ctx.Upgrade("name")`

## 功能

- `web.Context.Upgrade()`
- `web.Endpoint`
- `ws.Hook`
- `ws.Filter`
- `ws.Message`
- `ws.Command`
- `ctx.Reply()`
- `ctx.Broadcast()`
- `ctx.Groupcast()`
- `ctx.BindUser()`
- `ctx.PushUserResult()`
- 本地投递统计 `BroadcastResult / GroupcastResult`
- 协议导出 `/ws/export`
- 运行指标 `/ws/metrics`

## 运行

```bash
cd examples/ws-demo
go run .
```

默认监听 `http://127.0.0.1:8080/`。

页面里可以直接测试：

- 默认 ws 连接
- 自定义 endpoint 连接
- echo
- join / groupcast
- broadcast
- bind user / push user

其中：

- `/socket` 走 `ctx.Upgrade()`，命中 `ws` 模块注册的默认 Upgrade 接管器
- `/socket/custom` 走 `ctx.Upgrade("custom")`，命中 demo 自己注册的 `web.Endpoint`
