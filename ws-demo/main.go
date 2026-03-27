package main

import (
	"fmt"
	"time"

	_ "github.com/infrago/builtin"

	. "github.com/infrago/base"
	"github.com/infrago/infra"
	"github.com/infrago/web"
	"github.com/infrago/ws"
)

const demoPage = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Infrago WS Demo</title>
  <style>
    :root {
      --bg: #0b1020;
      --panel: #111934;
      --line: #273252;
      --text: #ebf0ff;
      --muted: #9aa7cc;
      --accent: #f59e0b;
      --accent2: #38bdf8;
      --button: #172347;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "IBM Plex Sans", "PingFang SC", sans-serif;
      color: var(--text);
      background:
        radial-gradient(circle at top left, rgba(56,189,248,.16), transparent 32%),
        radial-gradient(circle at top right, rgba(245,158,11,.18), transparent 28%),
        linear-gradient(160deg, #0a0f1d, #0f1832 48%, #0c1427);
      min-height: 100vh;
    }
    .wrap {
      max-width: 1040px;
      margin: 0 auto;
      padding: 40px 20px 72px;
    }
    .hero {
      margin-bottom: 24px;
      padding: 24px;
      border: 1px solid rgba(255,255,255,.08);
      border-radius: 24px;
      background: rgba(17,25,52,.76);
      backdrop-filter: blur(10px);
      box-shadow: 0 24px 80px rgba(0,0,0,.35);
    }
    h1 {
      margin: 0 0 8px;
      font-size: 32px;
      line-height: 1.1;
    }
    p {
      margin: 0;
      color: var(--muted);
      line-height: 1.6;
    }
    .grid {
      display: grid;
      grid-template-columns: 1.1fr .9fr;
      gap: 16px;
    }
    .card {
      padding: 20px;
      border-radius: 20px;
      border: 1px solid rgba(255,255,255,.08);
      background: rgba(17,25,52,.72);
      box-shadow: 0 16px 48px rgba(0,0,0,.28);
    }
    .row {
      display: flex;
      gap: 10px;
      margin-top: 12px;
      flex-wrap: wrap;
    }
    input {
      width: 100%;
      border: 1px solid var(--line);
      background: #0c1530;
      color: var(--text);
      border-radius: 14px;
      padding: 12px 14px;
      outline: none;
    }
    button {
      border: 0;
      background: var(--button);
      color: var(--text);
      padding: 12px 16px;
      border-radius: 14px;
      cursor: pointer;
      font-weight: 600;
    }
    button.primary {
      background: linear-gradient(135deg, #f59e0b, #ea580c);
      color: #111827;
    }
    button.secondary {
      background: linear-gradient(135deg, #38bdf8, #2563eb);
    }
    .status {
      margin-top: 12px;
      font-size: 14px;
      color: var(--muted);
    }
    pre {
      margin: 0;
      min-height: 420px;
      max-height: 560px;
      overflow: auto;
      border-radius: 16px;
      background: #09101f;
      border: 1px solid var(--line);
      padding: 16px;
      font-size: 13px;
      line-height: 1.5;
      color: #dbe6ff;
    }
    .label {
      margin-bottom: 10px;
      font-size: 12px;
      letter-spacing: .08em;
      text-transform: uppercase;
      color: var(--accent2);
    }
    @media (max-width: 900px) {
      .grid {
        grid-template-columns: 1fr;
      }
      pre {
        min-height: 280px;
      }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <section class="hero">
      <h1>WebSocket Demo</h1>
      <p>这个页面同时演示默认 <code>ctx.Upgrade()</code> 和自定义 <code>web.Endpoint</code> 接入。支持 echo、join、groupcast、broadcast。</p>
    </section>

    <section class="grid">
      <div class="card">
        <div class="label">Client</div>
        <div class="row">
          <button class="primary" onclick="connect()">Connect</button>
          <button onclick="connectCustom()">Connect Custom</button>
          <button onclick="disconnect()">Disconnect</button>
        </div>
        <div class="status" id="status">status: idle</div>

        <div class="row">
          <input id="echoText" placeholder="echo text" value="hello infrago">
        </div>
        <div class="row">
          <button class="secondary" onclick="sendEcho()">Send Echo</button>
        </div>

        <div class="row">
          <input id="groupId" placeholder="group id" value="room-demo">
        </div>
        <div class="row">
          <button onclick="joinGroup()">Join Group</button>
          <button onclick="groupcast()">Groupcast</button>
        </div>

        <div class="row">
          <input id="noticeText" placeholder="broadcast text" value="broadcast from browser">
        </div>
        <div class="row">
          <button onclick="broadcast()">Broadcast</button>
        </div>

        <div class="row">
          <input id="userId" placeholder="user id" value="user-demo">
        </div>
        <div class="row">
          <button onclick="bindUser()">Bind User</button>
          <button onclick="pushUser()">Push User</button>
        </div>
      </div>

      <div class="card">
        <div class="label">Events</div>
        <pre id="log"></pre>
      </div>
    </section>
  </div>

  <script>
    let ws;
    const logNode = document.getElementById("log");
    const statusNode = document.getElementById("status");

    function write(label, payload) {
      const line = "[" + new Date().toLocaleTimeString() + "] " + label + " " + JSON.stringify(payload);
      logNode.textContent = line + "\n" + logNode.textContent;
    }

    function connectTo(path) {
      if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
        return;
      }
      const schema = location.protocol === "https:" ? "wss://" : "ws://";
      ws = new WebSocket(schema + location.host + path);

      ws.onopen = () => {
        statusNode.textContent = "status: connected " + path;
        write("open", { ok: true, path: path });
      };
      ws.onclose = (event) => {
        statusNode.textContent = "status: closed";
        write("close", { code: event.code, reason: event.reason });
      };
      ws.onerror = () => {
        statusNode.textContent = "status: error";
        write("error", { ok: false });
      };
      ws.onmessage = (event) => {
        try {
          write("message", JSON.parse(event.data));
        } catch (e) {
          write("message", event.data);
        }
      };
    }

    function connect() {
      connectTo("/socket");
    }

    function connectCustom() {
      connectTo("/socket/custom");
    }

    function disconnect() {
      if (ws) {
        ws.close();
      }
    }

    function send(msg, args) {
      if (!ws || ws.readyState !== WebSocket.OPEN) {
        write("warn", { msg: "socket not connected" });
        return;
      }
      const payload = { name: msg, data: args };
      write("send", payload);
      ws.send(JSON.stringify(payload));
    }

    function sendEcho() {
      send("demo.echo", { text: document.getElementById("echoText").value });
    }

    function joinGroup() {
      send("demo.join", { group: document.getElementById("groupId").value });
    }

    function groupcast() {
      send("demo.groupcast", {
        group: document.getElementById("groupId").value,
        text: document.getElementById("noticeText").value
      });
    }

    function broadcast() {
      send("demo.broadcast", { text: document.getElementById("noticeText").value });
    }

    function bindUser() {
      send("demo.bind_user", { user: document.getElementById("userId").value });
    }

    function pushUser() {
      send("demo.push_user", {
        user: document.getElementById("userId").value,
        text: document.getElementById("noticeText").value
      });
    }
  </script>
</body>
</html>`

func main() {
	infra.Go()
}

func acceptOptions(ctx *web.Context, socket web.Socket) ws.AcceptOptions {
	return ws.AcceptOptions{
		Conn:       socket,
		Meta:       ctx.Meta,
		Name:       ctx.Name,
		Site:       ctx.Site,
		Host:       ctx.Host,
		Domain:     ctx.Domain,
		RootDomain: ctx.RootDomain,
		Path:       ctx.Path,
		Uri:        ctx.Uri,
		Setting:    ctx.Setting,
		Params:     ctx.Params,
		Query:      ctx.Query,
		Form:       ctx.Form,
		Value:      ctx.Value,
		Args:       ctx.Args,
		Locals:     ctx.Locals,
	}
}

func init() {
	infra.Register("custom", web.Endpoint{
		Name: "custom",
		Desc: "demo custom ws endpoint",
		Accept: func(ctx *web.Context, socket web.Socket) error {
			opts := acceptOptions(ctx, socket)
			locals := Map{}
			for key, value := range opts.Locals {
				locals[key] = value
			}
			locals["endpoint"] = "custom"
			opts.Locals = locals
			return ws.Accept(opts)
		},
	})

	infra.Register(".index", web.Router{
		Uri:  "/",
		Name: "ws demo home",
		Action: func(ctx *web.Context) {
			ctx.HTML(demoPage)
		},
	})

	infra.Register(".socket", web.Router{
		Uri:  "/socket",
		Name: "ws demo socket",
		Action: func(ctx *web.Context) {
			if err := ctx.Upgrade(); err != nil {
				ctx.Error(infra.Fail.With(err.Error()))
			}
		},
	})

	infra.Register(".socket.custom", web.Router{
		Uri:  "/socket/custom",
		Name: "ws demo custom socket",
		Action: func(ctx *web.Context) {
			if err := ctx.Upgrade("custom"); err != nil {
				ctx.Error(infra.Fail.With(err.Error()))
			}
		},
	})

	infra.Register(".ws.export", web.Router{
		Uri:  "/ws/export",
		Name: "ws demo export",
		Action: func(ctx *web.Context) {
			ctx.JSON(ws.Export())
		},
	})

	infra.Register(".ws.metrics", web.Router{
		Uri:  "/ws/metrics",
		Name: "ws demo metrics",
		Action: func(ctx *web.Context) {
			ctx.JSON(ws.Metrics())
		},
	})

	infra.Register("ws.access", ws.Hook{
		Name: "ws access",
		Open: func(ctx *ws.Context) {
			fmt.Printf("[ws] open sid=%s route=%s host=%s\n", ctx.Session.ID, ctx.Session.Name, ctx.Session.Host)
			endpoint, _ := ctx.Session.Locals["endpoint"].(string)
			if endpoint == "" {
				endpoint = "ws"
			}
			_ = ctx.Reply("demo.ready", Map{
				"sid":      ctx.Session.ID,
				"route":    ctx.Session.Name,
				"user":     ctx.Session.User,
				"endpoint": endpoint,
			})
		},
		Close: func(ctx *ws.Context) {
			fmt.Printf("[ws] close sid=%s route=%s result=%v\n", ctx.Session.ID, ctx.Session.Name, ctx.Result())
		},
		Receive: func(ctx *ws.Context) {
			fmt.Printf("[ws] receive sid=%s bytes=%d\n", ctx.Session.ID, len(ctx.Input))
		},
		Send: func(ctx *ws.Context) {
			fmt.Printf("[ws] send sid=%s msg=%s bytes=%d\n", ctx.Session.ID, ctx.Name, len(ctx.Output))
		},
	})

	infra.Register("ws.message.log", ws.Filter{
		Name: "ws message log",
		Message: func(ctx *ws.Context) {
			start := time.Now()
			fmt.Printf("[ws] message sid=%s msg=%s value=%v\n", ctx.Session.ID, ctx.Name, ctx.Value)
			ctx.Next()
			fmt.Printf("[ws] message done sid=%s msg=%s cost=%s\n", ctx.Session.ID, ctx.Name, time.Since(start))
		},
	})

	infra.Register("ws.error", ws.Handler{
		Name: "ws invalid handler",
		Invalid: func(ctx *ws.Context) {
			_ = ctx.Answer("demo.notice", nil, ctx.Result())
		},
		Error: func(ctx *ws.Context) {
			_ = ctx.Answer("demo.notice", nil, ctx.Result())
		},
	})

	infra.Register("demo.echo", ws.Message{
		Name: "echo",
		Action: func(ctx *ws.Context) {
			_ = ctx.Reply("demo.echoed", Map{
				"text": ctx.Value["text"],
				"sid":  ctx.Session.ID,
			})
		},
	})

	infra.Register("demo.join", ws.Message{
		Name: "join group",
		Action: func(ctx *ws.Context) {
			group, _ := ctx.Value["group"].(string)
			ctx.Join(group)
			_ = ctx.Reply("demo.notice", Map{
				"level": "info",
				"text":  fmt.Sprintf("joined group %s", group),
			})
		},
	})

	infra.Register("demo.groupcast", ws.Message{
		Name: "groupcast",
		Action: func(ctx *ws.Context) {
			group, _ := ctx.Value["group"].(string)
			result := ctx.GroupcastResult(group, "demo.notice", Map{
				"level": "group",
				"group": group,
				"text":  ctx.Value["text"],
				"sid":   ctx.Session.ID,
			})
			_ = ctx.Reply("demo.stats", Map{
				"op":      "groupcast",
				"group":   group,
				"hit":     result.Hit,
				"success": result.Success,
				"failed":  result.Failed,
				"error":   result.FirstError,
			})
		},
	})

	infra.Register("demo.broadcast", ws.Message{
		Name: "broadcast",
		Action: func(ctx *ws.Context) {
			result := ctx.BroadcastResult("demo.notice", Map{
				"level": "broadcast",
				"text":  ctx.Value["text"],
				"sid":   ctx.Session.ID,
			})
			_ = ctx.Reply("demo.stats", Map{
				"op":      "broadcast",
				"hit":     result.Hit,
				"success": result.Success,
				"failed":  result.Failed,
				"error":   result.FirstError,
			})
		},
	})

	infra.Register("demo.bind_user", ws.Message{
		Name: "bind user",
		Action: func(ctx *ws.Context) {
			user, _ := ctx.Value["user"].(string)
			ctx.BindUser(user)
			_ = ctx.Reply("demo.notice", Map{
				"level": "info",
				"text":  fmt.Sprintf("bound user %s", user),
				"user":  user,
			})
		},
	})

	infra.Register("demo.push_user", ws.Message{
		Name: "push user",
		Action: func(ctx *ws.Context) {
			user, _ := ctx.Value["user"].(string)
			result := ctx.PushUserResult(user, "demo.notice", Map{
				"level": "user",
				"text":  ctx.Value["text"],
				"user":  user,
				"sid":   ctx.Session.ID,
			})
			_ = ctx.Reply("demo.stats", Map{
				"op":      "push_user",
				"user":    user,
				"hit":     result.Hit,
				"success": result.Success,
				"failed":  result.Failed,
				"error":   result.FirstError,
			})
		},
	})

	infra.Register("demo.ready", ws.Command{Name: "ready"})
	infra.Register("demo.echoed", ws.Command{Name: "echoed"})
	infra.Register("demo.notice", ws.Command{Name: "notice"})
	infra.Register("demo.stats", ws.Command{Name: "stats"})
}
