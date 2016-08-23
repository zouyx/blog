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
