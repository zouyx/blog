package models

import (
	"testing"

	"github.com/astaxie/beego/orm"
)

func TestGetAllCategory(t *testing.T) {
	c, err := GetAllCategory()
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", c)
	}
}

func TestGetArticle(t *testing.T) {
	t.Log("query by id")
	a := &Article{Id_: 1}
	err := GetArticle(a, "")
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", a)
		t.Log("result tag:", a.Tags[0])
	}

	t.Log("query by name")
	b := &Article{Name: "joe"}
	err = GetArticle(b, "Name")
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", b)
	}
}

func TestGetArticleCount(t *testing.T) {
	i := GetArticleCount()
	t.Log("count:", i)
}

func TestGetAllTags(t *testing.T) {
	ts, err := GetAllTags()
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", ts)
	}
}

func TestGetArticlesByTag(t *testing.T) {
	a, _, err := GetArticlesByTag("2", 1, 1, "created_time")
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", a)
	}
}

func TestGetSubscribes(t *testing.T) {
	cond := orm.NewCondition()
	cond.And("status", false)
	a, _, err := GetSubscribes(cond, 0, 1, "")
	if err != nil {
		t.Error("error:", err)
	} else {
		t.Log("result:", a)
	}
}
