package root

import (
	"blog/common"
	"github.com/astaxie/beego"
)

type rootBaseRouter struct {
	beego.Controller
}

func (this *rootBaseRouter) Prepare() {
	if this.Ctx.Request.Method == "GET" {
		this.Data["SiteConfig"] = common.Webconfig
	}
	if this.Ctx.Request.URL.Path != "/root/login" {
		// sess_username, _ := this.GetSession("username").(string)
		sess_username := "joe"
		//如果未登录
		if sess_username == "" {
			this.Ctx.Redirect(302, "/root/login")
		} else {
			// this.SetSession("username", sess_username)
			this.Data["UserName"] = sess_username
		}
	}

}
