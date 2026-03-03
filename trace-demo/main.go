package main

import (
	"fmt"
	"time"

	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	"github.com/infrago/http"
	"github.com/infrago/infra"
	"github.com/infrago/log"
	_ "github.com/infrago/trace"
	_ "github.com/infrago/trace-file"
	_ "github.com/infrago/trace-greptime"
)

func main() {
	infra.Go()
}

func init() {

	infra.Register("index", http.Router{
		Uri: "/", Name: "index", Desc: "index",
		Action: func(ctx *http.Context) {
			ctx.Text("hello infra.")
		},
	})

	infra.Register("trace.child", infra.Service{
		Name: "子调用", Desc: "trace child service",
		Action: func(ctx *infra.Context) Map {
			ctx.Trace("搞飞机了这里")
			time.Sleep(10 * time.Millisecond)
			log.Debug("what")
			return Map{"ok": true, "at": time.Now().UnixMilli()}
		},
	})

	infra.Register(infra.START, infra.Trigger{
		Name: "Trace Demo",
		Desc: "emit trace spans on startup",
		Action: func(ctx *infra.Context) {
			data := ctx.Invoke("trace.child", Map{"from": "startup"})
			if res := ctx.Result(); res != nil && res.Fail() {
				return
			}

			span := ctx.Begin("开始")
			span.End()

			fmt.Println("trace demo done", data)
		},
	})
}
