package routers

import (
	"math"
	"strconv"

	"github.com/astaxie/beego/orm"

	"blog/common"
	"blog/models"
)

type NodeRouter struct {
	baseRouter
}

func (this *NodeRouter) Get() {
	limit := common.Webconfig.PageCount
	nodename := this.Ctx.Input.Param(":node")
	page, err := strconv.Atoi(this.Ctx.Input.Param(":page"))
	if err != nil {
		page = 1
	}
	cond := orm.NewCondition()
	cond.And("nname", nodename)

	articles, total, err := models.GetArticlesByNode(cond, (page-1)*limit, limit, "-created_time")

	if !this.CheckError(err) {
		return
	}

	if (page-1)*limit > total {
		this.Redirect("/prompt/404", 302)
		return
	}

	vars := make(map[string]interface{})

	totalpage := int(math.Ceil(float64(total) / float64(limit)))
	vars["CurrentCategory"] = &CurrentCategoryInfo{
		NName: nodename,
	}
	vars["Pager"] = common.GetPager("node/"+nodename, page, totalpage)
	data := MakeData(vars)

	this.Data["Data"] = data
	this.Data["Articles"] = articles
	// this.Data["json"] = articles
	// this.ServeJSON(true)
	this.Layout = "layout.html"
	this.TplName = "articles.html"
}
