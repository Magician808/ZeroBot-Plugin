package lolicon

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
	limit = rate.NewManager(time.Second*60, 5)
	m     = manager.New("setu\n 【发送/setu [关键词]】", &manager.Options{DisableOnDefault: false})
)

const recalltime = 60

func init() {
	engine := zero.New()

	single.New(
		single.WithKeyFn(func(ctx *zero.Ctx) interface{} {
			return ctx.Event.UserID
		}),
		single.WithPostFn(func(ctx *zero.Ctx) {
			//ctx.Send("您有操作正在执行，请稍后再试!")
			ctx.Send("[CQ:image,file=http://gchat.qpic.cn/gchatpic_new/0/0-2928072786-BA1AF4C420D2F485E7AA0C84C550C499/0?term=2]")

		}),
	).Apply(engine)

	_ = engine.OnCommandGroup([]string{"setu", "色图", "涩图", "搜图", "st"}).
		SetBlock(true).
		SetPriority(8).
		Handle(func(ctx *zero.Ctx) {
			var cmd extension.CommandModel
			err := ctx.Parse(&cmd)
			receivedmsgid := ctx.Event.MessageID
			if err != nil {
				ctx.Send(fmt.Sprintf("处理 %v 命令发生错误: %v", cmd.Command, err))
			}

			if cmd.Args == "" { // 未填写keyword,索取keyword
				ctx.Send(message.Message{message.Text("请输入setu关键词")})
				next := ctx.FutureEvent("message", ctx.CheckSession())
				recv, cancel := next.Repeat()
				for e := range recv {
					msg := e.Message.ExtractPlainText()
					if msg != "" {
						cmd.Args = msg
						receivedmsgid = e.MessageID
						cancel()
						continue
					}
					ctx.Send("关键词不合法oxo")
				}
			}
			zero.RangeBot(func(id int64, ctx2 *zero.Ctx) bool { // test the range bot function
				ctx.Send("准备图片ing,请稍等片刻(冷却:60s)")
				var pid, uid, title, author, urls string
				querylolicon(cmd.Args, &pid, &uid, &title, &author, &urls)
				if urls != "" {
					//msg_id := ctx.SendChain(message.Image(urls), message.Reply(receivedmsgid), message.Text(fmt.Sprintf("pid: %s\nuid: %s\ntitle: %s\nauthor: %s\nurl: %s", pid, uid, title, author, urls)))
					ctx.Send(message.ReplyWithMessage(receivedmsgid, message.Text(fmt.Sprintf("pid: %s\nuid: %s\ntitle: %s\nauthor: %s\nurl: %s", pid, uid, title, author, urls))))
					msg_id := ctx2.SendGroupMessage(ctx.Event.GroupID, message.Image(urls))
					/*rsp := ctx.CallAction("send_group_msg", zero.Params{
						"group_id": ctx.Event.GroupID,
						"message":  message.Image(querysetu(cmd.Args)),
					}).Data.Get("message_id")*/ //另一种实现方式
					time.Sleep(recalltime * time.Second)
					ctx2.DeleteMessage(msg_id)
				} else {
					ctx.Send("图片不存在或网络错误")
				}
				return true
			})
		})
	engine.UsePreHandler(m.Handler())

	engine.UsePreHandler(func(ctx *zero.Ctx) bool { // 限速器
		if !limit.Load(ctx.Event.UserID).Acquire() {
			ctx.Send("您的请求太快，请稍后重试0x0...")
			return false
		}
		return true
	})
}
