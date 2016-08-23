package routers

import "blog/models"

type ArticleRouter struct {
	baseRouter
}

func (this *ArticleRouter) Get() {
	name := this.Ctx.Input.Param(":article")

	article := &models.Article{Name: name}
	err := models.GetArticle(article,"Name")
	if !this.CheckError(err) {
		return
	}
	go article.UpdateViews()
	pre, next, _ := article.GetAroundArticle()
	vars := make(map[string]interface{})
	vars["PreArticle"] = pre
	vars["NextArticle"] = next
	vars["CurrentCategory"] = &CurrentCategoryInfo{
		ATitle: article.Title,
	}
	data := MakeData(vars)

	this.Data["Data"] = data
	this.Data["Article"] = article
	this.Data["SameTagArticls"] = article.GetSameTagArticles(5)
	this.Layout = "layout.html"
	this.TplName = "article.html"

	// this.Data["json"] = data
	// this.ServeJSON(true)
}
