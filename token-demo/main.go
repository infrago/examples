package main

import (
	"time"

	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	"github.com/infrago/http"
	"github.com/infrago/infra"
	_ "github.com/infrago/token"
	_ "github.com/infrago/token-memory"
)

func main() {
	infra.Go()
}

func init() {

	infra.Register("index", http.Router{
		Uri: "/", Name: "index", Desc: "index",
		Routing: http.Routing{
			"get": http.Router{
				Action: func(ctx *http.Context) {
					ctx.Text("index get")
				},
			},
			"post": http.Router{
				Action: func(ctx *http.Context) {
					ctx.Text("index post")
				},
			},
		},
	})

	infra.Register("all", http.Router{
		Uri: "/all", Name: "all", Desc: "all",
		Action: func(ctx *http.Context) {
			ctx.Text("index post")
		},
	})

	infra.Register("token.sign", http.Router{
		Uri:  "/token/sign",
		Name: "签发token",
		Desc: "issue one authed token",
		Action: func(ctx *http.Context) {
			uid := "u1001"
			if v, ok := ctx.Value["uid"].(string); ok && v != "" {
				uid = v
			}
			token := ctx.Sign(true, Map{
				"uid":   uid,
				"roles": []string{"user", "admin"},
			}, time.Hour)
			ctx.JSON(Map{
				"token":    token,
				"res":      ctx.Result().Error(),
				"token_id": ctx.TokenId(),
				"payload":  ctx.Payload(),
				"signed":   ctx.Signed(),
				"authed":   ctx.Authed(),
			})
		},
	})

	infra.Register("token.profile", http.Router{
		Uri:  "/token/profile",
		Name: "读取token",
		Desc: "requires valid signed token",
		Sign: true,
		Action: func(ctx *http.Context) {
			ctx.JSON(Map{
				"token_id": ctx.TokenId(),
				"payload":  ctx.Payload(),
				"signed":   ctx.Signed(),
				"authed":   ctx.Authed(),
			})
		},
	})

	infra.Register("token.auth", http.Router{
		Uri:  "/token/auth",
		Name: "认证路由",
		Desc: "requires authed=true",
		Auth: true,
		Action: func(ctx *http.Context) {
			ctx.JSON(Map{
				"ok":       true,
				"token_id": ctx.TokenId(),
				"payload":  ctx.Payload(),
			})
		},
	})

	infra.Register("token.revoke", http.Router{
		Uri:  "/token/revoke",
		Name: "吊销tokenId",
		Desc: "revoke current token id",
		Sign: true,
		Action: func(ctx *http.Context) {
			tid := ctx.TokenId()
			if tid != "" {
				_ = ctx.RevokeTokenID(tid, time.Now().Add(time.Hour).Unix())
			}
			ctx.JSON(Map{"ok": true, "token_id": tid, "revoked": true})
		},
	})
}
