package models

import (
	"html/template"
	"os"

	"blog/common"
	//"strings"
	//"log"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2/bson"
)

var (
	o         orm.Ormer
	DB_ENGINE = "INNODB"
	DRIVER    = "mysql"
	DEBUG     = true
)

func init() {
	orm.RegisterDriver(DRIVER, orm.DRMySQL)

	orm.RegisterDataBase("default", DRIVER, "blog:123456@tcp(10.0.12.19:3306)/blogdb?charset=utf8")
	orm.RegisterModel(new(Article),
		new(Comment),
		new(Category),
		new(TagWrapper),
		new(Subscription),
		new(Node))
	orm.RunSyncdb("default", false, true)

	o = orm.NewOrm()

	if DEBUG {
		orm.Debug = DEBUG
		// var w io.Writer

		orm.DebugLog = orm.NewLog(os.Stderr)
	}
}

type Article struct {
	Id_            int64 `orm:"column(id);pk;auto"`
	CName          string
	NName          string
	Name           string
	Author         string
	Title          string
	Text           template.HTML `orm:"null"`
	Tags           []string      `orm:"-"`
	FeaturedPicURL string        `orm:"column(featured_pic_url)"`
	Summary        template.HTML `orm:"null"`
	Views          int           `orm:"null"`
	Comments       []*Comment    `orm:"reverse(many);null"`
	IsThumbnail    bool          `orm:"null"`
	CreatedTime    time.Time
	ModifiedTime   time.Time
}

func (this *Article) TableEngine() string {
	return DB_ENGINE
}

func (article *Article) SetSummary() {
	if article.IsThumbnail {
		article.Summary = article.Text
	} else {
		strs := strings.Split(string(article.Text), "<!--more-->")
		//beego.Error(strs[0])
		n := len(strs)
		if n > 0 {
			article.Summary = template.HTML(strs[0])
		}
	}
}

func (article *Article) GetFirstParagraph() *template.HTML {
	rx := regexp.MustCompile(`<p>(.*)</p>`)
	p := rx.FindStringSubmatch(string(article.Text))
	//beego.Error(p)
	n := len(p)
	if n > 1 {
		rep := template.HTML(p[1] + "...")
		return &rep
	}
	return nil
}

func (article *Article) GetCategory() *Category {
	// c := DB.C("category")
	// var category Category
	// c.Find(bson.M{"name": article.CName}).One(&category)
	// return &category
	var category Category
	for _, v := range Categories {
		if v.Name == article.CName {
			category = v
			break
		}
	}
	return &category
}

func (article *Article) GetNode() *Node {
	var node *Node
	for _, v := range Categories {
		if v.Name == article.CName {
			for _, va := range v.Nodes {
				if va.Name == article.NName {
					node = va
					break
				}
			}
			break
		}
	}

	return node
}

func (article *Article) GetTags() *[]TagWrapper {
	return article.GetSelfTags()
}

func (article *Article) CreatArticle() error {
	// article.Id_ = bson.NewObjectId()
	_, err := o.Insert(article)
	// go setTags(&article.Tags, article.Id_)
	return err
}

func (article *Article) UpdateArticle() error {
	c := DB.C("article")
	err := c.UpdateId(article.Id_, article)
	// go setTags(&article.Tags, article.Id_)

	return err
}

func (article *Article) GetCommentCount() int {
	return 1
}

func (article *Article) GetAroundArticle() (*Article, *Article, error) {
	// c := DB.C("article")
	var preresult, nextresult Article
	// err := c.Find(&bson.M{"createdtime": &bson.M{"$lt": article.CreatedTime}}).Sort("-createdtime").Limit(1).One(&preresult)

	// err = c.Find(&bson.M{"createdtime": &bson.M{"$gt": article.CreatedTime}}).Sort("createdtime").Limit(1).One(&nextresult)
	err := o.Raw("SELECT  * FROM article WHERE created_time<? Order by created_time desc limit 1 ", article.CreatedTime.String()).QueryRow(&preresult)
	err = o.Raw("SELECT  * FROM article WHERE created_time>? Order by created_time limit 1 ", article.CreatedTime.String()).QueryRow(&nextresult)

	return &preresult, &nextresult, err
}

func (article *Article) GetSameTagArticles(limit int) (articles []Article) {
	ids := make([]int64, 0)
	// for _, v := range Tags {
	// 	for _, tag := range article.Tags {
	// 		if tag == v.Title || tag == v.Name {
	// 			for _, va := range v.ArticleIds {
	// 				if va != article.Id_ {
	// 					ids = append(ids, va)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	d := DB.C("article")
	d.Find(&bson.M{"_id": &bson.M{"$in": ids}}).Limit(limit).All(&articles)
	return
}

func (article *Article) GetSelfTags() *[]TagWrapper {
	var tags []TagWrapper
	for _, v := range Tags {
		for _, va := range v.ArticleIds {
			if va != article.Id_ {
				tags = append(tags, v)
			}
		}
	}
	return &tags
}

func (article *Article) HasFeaturedPic() bool {
	if len(article.FeaturedPicURL) == 0 {
		return false
	}
	return true
}

func (article *Article) HasSummary() bool {
	if len(article.Summary) == 0 {
		return false
	}
	return true
}

func (article *Article) UpdateViews() {
	article.Views++
	article.UpdateArticle()
}

type CommentIndexItem struct {
	Name         string
	CommentNames []string
}

type CommentMetadata struct {
	Name        string
	Author      string
	ArticleName string
	UAgent      string `orm:"null"`
	URL         string `orm:"column(url);null"`
	IP          string `orm:"column(ip);null"`
	Email       string `orm:"null"`
	EmailHash   string `orm:"null"`
	CreatedTime int64
}

func (m *Comment) BuildFromJson(json interface{}) {
	var jsonMap map[string]interface{}
	jsonMap = json.(map[string]interface{})
	for k, v := range jsonMap {
		switch vv := v.(type) {
		case string:
			str := vv
			switch k {
			case "Name":
				m.Name = str
			case "Author":
				m.Author = str
			case "URL":
				m.URL = str
			case "IP":
				m.IP = str
			case "Email":
				m.Email = str
			case "EmailHash":
				m.EmailHash = str
			case "UAgent":
				m.UAgent = str
			case "ArticleName":
				m.ArticleName = str
			}
		case float64:
			if k == "CreatedTime" {
				m.CreatedTime = int64(vv)
			}
		default:
		}
	}
}

func (meta *Comment) CreatedTimeHumanReading() string {
	return common.TimeHumanReading(meta.CreatedTime)
}

type Comment struct {
	Id      int64    `orm:"column(id);pk;auto"`
	Article *Article `orm:"rel(fk);null"`
	CommentMetadata
	Text template.HTML `orm:"null"`
}

func (this *Comment) TableEngine() string {
	return DB_ENGINE
}

func (this *Comment) CreateComment() error {
	_, err := o.Insert(this)
	return err
}

func (this *Comment) TableName() string {
	return "article_comment"
}

type Node struct {
	Id          int64     `orm:"column(id);pk;auto"`
	Category    *Category `orm:"rel(fk);null"`
	Name        string
	Title       string
	Content     string
	CreatedTime time.Time `orm:"null"`
	UpdatedTime time.Time `orm:"null"`
	Views       int64     `orm:"null"`
	ArticleTime time.Time `orm:"null"`
}

func (this *Node) CreateNode() error {
	_, err := o.Insert(this)
	return err
}

func (this *Node) TableEngine() string {
	return DB_ENGINE
}

func (this *Node) TableName() string {
	return "node"
}

func (node *Node) GetAllArticles(offset int, limit int, sort string) (*[]Article, int, error) {
	c := DB.C("article")

	var article []Article
	q := bson.M{"nname": node.Name}
	query := c.Find(q).Skip(offset).Limit(limit)
	if sort != "" {
		query = query.Sort(sort)
	}
	err := query.All(&article)
	total, _ := c.Find(q).Count()
	return &article, total, err
}

func (node *Node) GetArticleCount() (int, error) {
	qs := o.QueryTable("article")
	num, err := qs.Filter("n_name", node.Name).Count()
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

// func (node *Node) GetCategory() (*Category, error) {
// 	c := DB.C("category")
// 	var category Category
// 	err := c.Find(bson.M{"_id": node.CId_}).One(&category)
// 	return &category, err
// }

// func (node *Node) CreatNode() error {
// 	//node.Id_ = bson.NewObjectId()
// 	c := DB.C("node")
// 	err := c.Insert(node)
// 	return err
// }

// func (node *Node) UpdateNode() error {
// 	c := DB.C("node")
// 	err := c.UpdateId(node.Id_, node)
// 	return err
// }

type Category struct {
	Id_         int64 `orm:"column(id);pk;auto"`
	Name        string
	Title       string
	Content     string `orm:"null"`
	CreatedTime time.Time
	UpdatedTime time.Time
	Views       int       `orm:"null"`
	NodeTime    time.Time `orm:"null"`
	Nodes       []*Node   `orm:"reverse(many);null"`
}

func (this *Category) TableEngine() string {
	return DB_ENGINE
}

func (this *Category) CreatCategory() error {
	_, err := o.Insert(this)
	SetAppCategories()
	return err
}

func (this *Category) UpdateCategory() error {
	_, err := o.Update(this)
	SetAppCategories()
	return err
}

// func (category *Category) GetAllNodes() *[]Node {
// 	c := DB.C("node")
// 	var nodes []Node
// 	c.Find(&bson.M{"_cid": category.Id_}).All(&nodes)
// 	return &nodes
// }

// func (category *Category) SetNodeId(nid bson.ObjectId) {
// 	if category.NodeIds != nil {
// 		category.NodeIds = append(category.NodeIds, nid)
// 		removeDuplicate(&category.NodeIds)
// 	} else {
// 		category.NodeIds = []bson.ObjectId{nid}
// 	}
// 	category.NodeCount = len(category.NodeIds)
// }

type TagWrapper struct {
	Id_          int64 `orm:"column(id);pk;auto"`
	Name         string
	Title        string
	Count        int
	CreatedTime  time.Time
	ModifiedTime time.Time
	ArticleIds   []int64 `orm:"-"`
}

func (this *TagWrapper) TableEngine() string {
	return DB_ENGINE
}

func (tag *TagWrapper) SetTag() error {
	// c := DB.C("tags")
	var err error
	flag := false
	for _, v := range Tags {
		if tag.Name == v.Name {
			v.ArticleIds = append(v.ArticleIds, tag.ArticleIds...)
			removeDuplicate(&v.ArticleIds)
			v.Count = len(v.ArticleIds)
			v.ModifiedTime = time.Now()
			_, err = o.Update(v)
			flag = true
			break
		}
	}

	if !flag {
		// tag.Id_ = bson.NewObjectId()
		tag.CreatedTime = time.Now()
		tag.ModifiedTime = time.Now()
		Tags = append(Tags, *tag)
		_, err = o.Insert(tag)
	}

	SetAppTags()
	return err
}

type Subscription struct {
	Id_    int64 `orm:"column(id);pk;auto"`
	Email  string
	Uid    string
	Status bool
}

func (this *Subscription) TableEngine() string {
	return DB_ENGINE
}

func (this *Subscription) Set() error {
	// c := DB.C("subscription")
	sub := &Subscription{
		Email:  this.Email,
		Status: this.Status,
		Uid:    this.Uid,
	}
	created, id, err := o.ReadOrCreate(sub, "Email")
	if err == nil {
		if !created {
			sub.Id_ = id
			sub.Email = this.Email
			sub.Status = this.Status
			sub.Uid = this.Uid
			_, err = o.Update(sub)
		}
	}
	return err
}

func (this *Subscription) UpdateState() error {
	_, err := o.QueryTable("subscription").Filter("Uid", this.Uid).Update(orm.Params{
		"status": this.Status,
	})
	return err
}
