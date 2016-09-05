package models

//"strings"

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var (
	Categories []Category   //常驻内存
	Tags       []TagWrapper //常驻内存
)

func InitDb() {
	SetAppCategories()
	SetAppTags()
}

func SetAppCategories() {
	Categories, _ = GetAllCategory()
}

func SetAppTags() {
	tags, _ := GetAllTags()
	Tags = *tags
}

func GetArticles(condition *orm.Condition, offset int, limit int, sort string) (*[]Article, int, error) {
	var articles []Article
	//下面为新增
	qs := o.QueryTable("article")
	qs = qs.SetCond(condition).Offset(offset).Limit(limit)
	if sort != "" {
		qs = qs.OrderBy(sort)
	}

	_, err := qs.All(&articles)

	total, _ := qs.SetCond(condition).Count()
	if err != nil {
		return nil, 0, err
	}
	return &articles, int(total), nil
}

func GetArticlesByTag(tagname string, offset int, limit int, sort string) (*[]Article, int, error) {
	if len(Tags) == 0 {
		articles := make([]Article, 0)
		return &articles, 0, nil
	}
	var tag TagWrapper
	for _, v := range Tags {
		if tagname == v.Name {
			tag = v
			break
		}
	}
	_, err := o.LoadRelated(&tag, "Articles")
	if err != nil {
		return nil, 0, err
	}

	if len(tag.Articles) == 0 {
		articles := make([]Article, 0)
		return &articles, 0, nil
	}
	articleIds := make([]int64, 0, len(tag.Articles))
	for _, article := range tag.Articles {
		articleIds = append(articleIds, article.Id_)
	}
	cond := orm.NewCondition()
	cond.And("id__in", articleIds)
	return GetArticles(cond, offset, limit, sort)
}

func GetArticlesByNode(condition *orm.Condition, offset int, limit int, sort string) (*[]Article, int, error) {
	return GetArticles(condition, offset, limit, sort)
}

func GetArticleCount() int {
	cnt, err := o.QueryTable("article").Count()
	if err != nil {
		return 0
	}
	return int(cnt)
}

func GetArticle(article *Article, column string) error {
	var err error
	if column == "" {
		err = o.Read(article)
	} else {
		err = o.Read(article, column)
	}
	if article.Id_ > 0 {
		_, err = o.LoadRelated(article, "Tags")
	}
	return err
}

// func DeleteArticles(condition *orm.Condition) (*mgo.ChangeInfo, error) {
// 	// c := DB.C("article")
// 	qs := o.QueryTable("article")
// 	num, err := qs.SetCond(condition).Delete()
// 	return &mgo.ChangeInfo{}, nil
// }

func DeleteArticle(id int64) error {
	_, err := o.Delete(&Article{Id_: id})
	return err
}

// func GetTags(condition *orm.Condition, offset int, limit int, sort string) (*[]TagWrapper, int, error) {
// 	var tags []TagWrapper
// 	qs := o.QueryTable("tag_wrapper")
// 	qs = qs.SetCond(condition).Offset(offset).Limit(limit)
// 	if sort != "" {
// 		qs = qs.OrderBy(sort)
// 	}

// 	_, err := qs.All(&tags)

// 	total, _ := qs.SetCond(condition).Count()
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return &tags, int(total), nil
// }

func GetAllTags() (*[]TagWrapper, error) {
	var tags []TagWrapper
	o.QueryTable("tag_wrapper").All(&tags)
	return &tags, nil
}

func GetTagCount() int {
	return len(Tags)
}

func GetAllCategory() ([]Category, error) {
	qb, _ := orm.NewQueryBuilder(DRIVER)

	qb.Select(" * ").
		From("category").
		OrderBy("created_time").Desc()
	var categories []Category
	_, err := o.Raw(qb.String()).QueryRows(&categories)
	if err != nil {
		return nil, err
	}
	//获取关系表信息
	for _, category := range categories {
		category.LoadNodes()
	}
	return categories, nil
}

func GetCategoryById(id int64) Category {
	var category Category
	for _, v := range Categories {
		if v.Id_ == id {
			category = v
			break
		}
	}
	return category
}
func GetCategoryNodeName(nname string) Category {
	var category Category
	for _, v := range Categories {
		flag := false
		for _, va := range v.Nodes {
			if va.Name == nname {
				flag = true
				break
			}
		}
		if flag {
			category = v
		}
	}
	return category
}

// func DeleteCategory(condition *bson.M) error {
// 	// c := DB.C("category")
// 	// err := c.Remove(condition)
// 	// SetAppCategories()
// 	return nil
// }

func GetSubscribes(condition *orm.Condition, offset int, limit int, sort string) (*[]Subscription, int, error) {
	var subs []Subscription
	qs := o.QueryTable("subscription")
	qs = qs.SetCond(condition).Offset(offset).Limit(limit)
	if sort != "" {
		qs = qs.OrderBy(sort)
	}

	_, err := qs.All(&subs)

	total, _ := qs.SetCond(condition).Count()
	if err != nil {
		return nil, 0, err
	}
	return &subs, int(total), nil
}

func removeDuplicate(slis []*Article) []*Article {
	newArticles := make([]*Article, 0)
	for i, val := range slis {
		found := false
		for j, jval := range slis {
			if val.Id_ == jval.Id_ && i != j {
				found = true
				break
			}
		}

		if !found {
			newArticles = append(newArticles, val)
		}
	}

	return newArticles
}

func setTags(tags []*TagWrapper, article *Article) {
	var tag *TagWrapper
	for _, v := range tags {
		if v.Id_ == 0 {
			tag = &TagWrapper{
				Name:     v.Name,
				Title:    v.Title,
				Count:    1,
				Articles: []*Article{article},
			}
		} else {
			tag = &TagWrapper{
				Id_: v.Id_,
			}
			err := o.Read(tag)
			if err != nil {
				beego.Error("Read tag error:", err)
			}
			tag.Articles = []*Article{article}
		}
		tag.SetTag()
	}
}
