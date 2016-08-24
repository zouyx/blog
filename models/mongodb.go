package models

//"strings"

import (
	"github.com/astaxie/beego/orm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	DB         *mgo.Database
	Categories []Category   //常驻内存
	Tags       []TagWrapper //常驻内存
)

func InitDb() {

	// conn := common.Webconfig.Dbconn
	// if conn == "" {
	// 	beego.Error("数据库地址还没有配置,请到config内配置db字段.")
	// 	os.Exit(1)
	// }

	// session, err := mgo.Dial(conn)
	// if err != nil {
	// 	beego.Error("MongoDB连接失败:", err.Error())
	// 	os.Exit(1)
	// }

	// session.SetMode(mgo.Monotonic, true)

	// DB = session.DB("messageblog")
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

func GetArticles(condition *bson.M, offset int, limit int, sort string) (*[]Article, int, error) {
	// c := DB.C("article")
	var article []Article
	// query := c.Find(condition).Skip(offset).Limit(limit)
	// if sort != "" {
	// 	query = query.Sort(sort)
	// }
	// err := query.All(&article)
	// total, _ := c.Find(condition).Count()
	article = make([]Article, 1, 1)
	return &article, 0, nil
}

func GetArticlesByTag(tagname string, offset int, limit int, sort string) (*[]Article, int, error) {
	var tag TagWrapper
	for _, v := range Tags {
		if tagname == v.Name {
			tag = v
			break
		}
	}
	return GetArticles(&bson.M{"_id": bson.M{"$in": tag.ArticleIds}}, offset, limit, sort)
}

func GetArticlesByNode(condition *bson.M, offset int, limit int, sort string) (*[]Article, int, error) {
	return GetArticles(condition, offset, limit, sort)
}

func GetArticleCount() int {
	// c := DB.C("article")
	// total, _ := c.Count()
	cnt, err := o.QueryTable("article").Count()
	if err != nil {
		return 0
	}
	return int(cnt)
}

func GetArticle(article *Article, column string) error {
	if column == "" {
		return o.Read(article)
	}
	return o.Read(article, column)
}

func DeleteArticles(condition *bson.M) (*mgo.ChangeInfo, error) {
	// c := DB.C("article")
	return &mgo.ChangeInfo{}, nil
}

func DeleteArticle(id int64) error {
	_, err := o.Delete(&Article{Id_: id})
	return err
}

func GetTags(condition *bson.M, offset int, limit int, sort string) (*[]TagWrapper, int, error) {
	// c := DB.C("tags")
	// var tags []TagWrapper
	// query := c.Find(condition).Skip(offset).Limit(limit)
	// if sort != "" {
	// 	query = query.Sort(sort)
	// }
	// err := query.All(&tags)
	// total, _ := c.Find(condition).Count()

	var tags []TagWrapper
	tags = make([]TagWrapper, 1, 1)

	return &tags, 0, nil
}

func GetAllTags() (*[]TagWrapper, error) {
	// c := DB.C("tags")
	// var tags []TagWrapper
	// err := c.Find(&bson.M{}).All(&tags)
	var tags []TagWrapper
	tags = make([]TagWrapper, 1, 1)
	return &tags, nil
}

func GetTagCount() int {
	return len(Tags)
}

func GetAllCategory() ([]Category, error) {
	// c := DB.C("category")
	// var categories []Category
	// err := c.Find(bson.M{}).Sort("createdtime").All(&categories)
	qb, _ := orm.NewQueryBuilder(DRIVER)

	qb.Select(" * ").
		From("category").
		OrderBy("created_time").Desc()
	var categories []Category
	_, err := o.Raw(qb.String()).QueryRows(&categories)
	if err != nil {
		return nil, err
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

func DeleteCategory(condition *bson.M) error {
	// c := DB.C("category")
	// err := c.Remove(condition)
	// SetAppCategories()
	return nil
}

func GetSubscribes(condition *bson.M, offset int, limit int, sort string) (*[]Subscription, int, error) {
	// c := DB.C("subscription")
	// var subs []Subscription
	// query := c.Find(condition).Skip(offset).Limit(limit)
	// if sort != "" {
	// 	query = query.Sort(sort)
	// }
	// err := query.All(&subs)
	// total, _ := c.Find(condition).Count()
	var subs []Subscription
	subs = make([]Subscription, 1, 1)
	return &subs, 0, nil
}

func removeDuplicate(slis *[]int64) {
	found := make(map[int64]bool)
	j := 0
	for i, val := range *slis {
		if _, ok := found[val]; !ok {
			found[val] = true
			(*slis)[j] = (*slis)[i]
			j++
		}
	}
	*slis = (*slis)[:j]
}

func setTags(tags *[]string, aid int64) {
	for _, v := range *tags {
		tag := &TagWrapper{
			Name:       v,
			Title:      v,
			Count:      1,
			ArticleIds: []int64{aid},
		}
		tag.SetTag()
	}
}
