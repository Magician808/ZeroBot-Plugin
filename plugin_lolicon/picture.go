package lolicon

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {

	zero.OnFullMatch("随机美图", zero.OnlyToMe).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send(message.Image("http://iw233.cn/api/Random.php"))
		})
}
