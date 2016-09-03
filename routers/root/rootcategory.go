package root

import (
	"strconv"
	"time"

	"blog/models"
)

type RootCategoryRouter struct {
	rootBaseRouter
}

func (this *RootCategoryRouter) Get() {
	id, err := this.GetInt64("id")

	if err != nil && id > 0 {
		for _, v := range models.Categories {
			if v.Id_ == id {
				this.Data["json"] = v
				break
			}
		}
		this.ServeJSON(true)
	} else {
		this.Data["Category"] = models.Categories
		this.Data["currentitem"] = "category"
		this.Layout = "root/layout.html"
		this.TplName = "root/category.html"
	}
}

func (this *RootCategoryRouter) Post() {
	id, err := this.GetInt64("id")

	if len(this.Input()) == 1 { //删除操作
		// models.DeleteCategory(&time.M{"_id": time.ObjectIdHex(id)})
		this.Data["json"] = true
		this.ServeJSON(true)
	} else {
		name := this.GetString("name")
		title := this.GetString("title")
		content := this.GetString("content")
		if name == "" {
			name = strconv.Itoa(int(time.Now().UnixNano()))
		}
		if err != nil && id > 0 {
			for _, v := range models.Categories {
				if v.Id_ == id {
					v.Name = name
					v.Title = title
					v.Content = content
					v.UpdatedTime = time.Now()
					v.UpdateCategory()
					break
				}
			}
		} else {
			cat := models.Category{
				Name:        name,
				Title:       title,
				Content:     content,
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
				NodeTime:    time.Now(),
			}
			cat.CreatCategory()
		}

		this.Redirect("/root/category", 302)
	}
}
