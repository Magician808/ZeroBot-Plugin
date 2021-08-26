package manager

import (
	"fmt"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/example/manager"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := zero.New()
	var m = manager.New("Repeater/n 【发送/[复读前缀] 开启复读】\n 【发送/[取消前缀] 关闭复读】", &manager.Options{DisableOnDefault: false})
	zero.OnCommandGroup([]string{"repeat start", "开始复读"}).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
			stop := zero.NewFutureEvent("message", 8, true,
				zero.CommandRule("repeat stop"), // 关闭复读指令
				ctx.CheckSession()).             // 只有开启者可以关闭复读模式
				Next()                           // 关闭需要一次

			echo, cancel := ctx.FutureEvent("message",
				ctx.CheckSession()). // 只复读开启复读模式的人的消息
				Repeat()             // 不断监听复读
			msg_id := ctx.Send("已开启复读模式!")
			time.Sleep(2 * time.Second)
			ctx.DeleteMessage(msg_id)
			ctx.Send(fmt.Sprintf("[CQ:poke,qq=%d]", ctx.Event.SelfID))
			for {
				select {
				case e := <-echo: // 接收到需要复读的消息
					ctx.Send(message.UnescapeCQText(e.RawMessage))
				case <-stop: // 收到关闭复读指令
					cancel() // 取消复读监听
					return   // 返回
				}
			}
		})
	engine.UsePreHandler(m.Handler())
}
