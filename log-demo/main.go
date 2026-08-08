package main

import (
	"sync"
	"time"

	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	"github.com/infrago/http"
	"github.com/infrago/infra"
	"github.com/infrago/log"
	_ "github.com/infrago/log-file"
	_ "github.com/infrago/log-greptime"
)

func main() {
	infra.Go()
}

func init() {

	infra.Register("index", http.Router{
		Uri: "/", Name: "index", Desc: "index",
		Action: func(ctx *http.Context) {
			log.Debugw("index", Map{"host": ctx.Host})
			ctx.Text("hello infra.")
		},
	})

	infra.Register("stats", http.Router{
		Uri:  "/stats",
		Name: "日志统计",
		Desc: "log stats",
		Action: func(ctx *http.Context) {
			ctx.JSON(log.Stats())
		},
	})

	infra.Register(infra.START, infra.Trigger{
		Name: "Log Demo",
		Desc: "emit many logs to test async batch and overflow strategy",
		Action: func(ctx *infra.Context) {
			const workers = 16
			const perWorker = 5000

			start := time.Now()
			var wg sync.WaitGroup
			wg.Add(workers)

			for w := 0; w < workers; w++ {
				worker := w
				go func() {
					defer wg.Done()
					for i := 0; i < perWorker; i++ {
						if i%1000 == 0 {
							log.Warningw("burst", Map{"worker": worker, "seq": i})
						} else {
							log.Infof("burst worker=%d seq=%d", worker, i)
						}
					}
				}()
			}

			wg.Wait()
			took := time.Since(start)
			log.Noticew("log-demo done", Map{"workers": workers, "per_worker": perWorker, "total": workers * perWorker, "cost": took.String()})
		},
	})
}
