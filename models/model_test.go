package models

import (
	"testing"
	"time"
)

func TestCreateArticle(t *testing.T) {
	a := &Article{}
	a.Author = "joe"
	a.CName = "joe"
	a.NName = "joe"
	a.Name = "joe"
	a.CreatedTime = time.Now()
	a.ModifiedTime = time.Now()

	// comments := make([]*Comment, 0)

	// comments = append(comments, c)
	// a.Comments = comments

	err := a.CreatArticle()
	t.Log(err)

	c := &Comment{}
	c.Author = "mic"
	c.Name = "mic"
	c.ArticleName = "mic"
	c.CreatedTime = time.Now().Unix()
	c.Article = a
	e := c.CreateComment()
	t.Log(e)
}

// func TestCreateComment(t *testing.T) {
// 	c := &Comment{}
// 	c.Author = "mic"
// 	c.Name = "mic"
// 	c.ArticleName = "mic"
// 	c.CreatedTime = time.Now().Unix()
// 	err := c.CreateComment()
// 	t.Error(err)
// }

func TestCreateCategory(t *testing.T) {
	a := &Category{}
	a.Name = "joe"
	a.Title = "joe"
	a.CreatedTime = time.Now()
	a.UpdatedTime = time.Now()
	err := a.CreatCategory()
	t.Log(err)

	n := &Node{}
	n.Name = "joe"
	n.Content = "joe"
	n.Title = "joe"
	n.Category = a
	err = n.CreateNode()
	t.Log(err)
}