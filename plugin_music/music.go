package music

import (
	"fmt"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/example/manager"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	limit2 = rate.NewManager(time.Second*60, 2)
	m     = manager.New("Music\n 【发送/music [音乐名]】", &manager.Options{DisableOnDefault: false})
)

func init() {
	engine := zero.New()

	single.New(
		single.WithKeyFn(func(ctx *zero.Ctx) interface{} {
			return ctx.Event.UserID
		}),
		single.WithPostFn(func(ctx *zero.Ctx) {
			ctx.Send("您有操作正在执行，请稍后再试!")
		}),
	).Apply(engine)

	_ = engine.OnCommandGroup([]string{"music", "点歌"}).
		SetBlock(true).
		SetPriority(8).
		Handle(func(ctx *zero.Ctx) {
			var cmd extension.CommandModel
			err := ctx.Parse(&cmd)
			if err != nil {
				ctx.Send(fmt.Sprintf("处理 %v 命令发生错误: %v", cmd.Command, err))
			}

			if cmd.Args == "" { // 未填写歌曲名,索取歌曲名
				ctx.Send(message.Message{message.Text("请输入要点的歌曲!")})
				next := ctx.FutureEvent("message", ctx.CheckSession())
				recv, cancel := next.Repeat()
				for e := range recv {
					msg := e.Message.ExtractPlainText()
					if msg != "" {
						cmd.Args = msg
						cancel()
						continue
					}
					ctx.Send("歌曲名不合法oxo")
				}
			}
			zero.RangeBot(func(id int64, ctx2 *zero.Ctx) bool { // test the range bot function
				ctx2.SendGroupMessage(ctx.Event.GroupID, message.Music("163", queryNeteaseMusic(cmd.Args)))
				return true
			})
			// ctx.Send(message.Music("163", queryNeteaseMusic(cmd.Args)))
		})
	engine.UsePreHandler(m.Handler())

	engine.UsePreHandler(func(ctx *zero.Ctx) bool { // 限速器
		if !limit2.Load(ctx.Event.UserID).Acquire() {
			ctx.Send("您的请求太快，请稍后重试0x0...")
			return false
		}
		return true
	})
}
