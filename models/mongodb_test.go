package models

import "testing"

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
