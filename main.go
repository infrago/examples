package main

import (
	"fmt"
	"time"

	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"

	"github.com/infrago/cron"
	"github.com/infrago/log"
	"github.com/infrago/mutex"

	"github.com/infrago/http"
	"github.com/infrago/infra"

	_ "github.com/infrago/cron-postgres"
)

func main() {
	infra.Go()
}

func init() {

	infra.Register("ssss.test", infra.Service{
		Name: "test", Desc: "test",
		Action: func(ctx *infra.Context) (Map, Res) {
			log.Debug("ssss.test", time.Now())
			return nil, infra.OK
		},
	})

	infra.Register("cron.test", infra.Method{
		Name: "test", Desc: "test",
		Action: func(ctx *infra.Context) (Map, Res) {
			log.Debug("cron.test", time.Now())
			return nil, infra.OK
		},
	})

	infra.Register("test", cron.Job{
		Schedule: "*/10 * * * * *", Target: "cron.test",
	})

	infra.Register("index", http.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *http.Context) {
			jobs := cron.ListJobs()
			count, logs := cron.ListLogs("test", 0, 10)
			ctx.JSON(Map{
				"count": count, "logs": logs,
				"jobs": jobs,
			})
		},
	})

	// infra.Register("www.index", http.Router{
	// 	Uri: "/", Name: "首页", Desc: "首页",
	// 	Action: func(ctx *http.Context) {
	// 		cache.Write("key", Map{"msg": "msg from cache."}, time.Second*10)
	// 		ctx.Text("hello world.")
	// 	},
	// })
	// infra.Register("www.json", http.Router{
	// 	Uri: "/json", Name: "JSON", Desc: "JSON",
	// 	Action: func(ctx *http.Context) {
	// 		data, _ := cache.Read("key")
	// 		ctx.Answer(nil, Map{
	// 			"msg":   "hello world.",
	// 			"cache": data,
	// 		})
	// 	},
	// })

	infra.Register(infra.START, infra.Trigger{
		Name: "启动", Desc: "启动",
		Action: func(ctx *infra.Context) {
			data := ctx.Invoke("ssss.test", Map{"msg": "msg from examples."})
			res := ctx.Result()

			fmt.Println("ssss", res, data)

			_, err := mutex.Lock("test", time.Minute)
			fmt.Println("lock", err)
		},
	})

}
