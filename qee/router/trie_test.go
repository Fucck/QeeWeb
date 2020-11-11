package router

import (
	"testing"
)

func TestTrieRouter_RegisteredHandler_absParse(t *testing.T) {
	pattern := "/hello/world"
	router := NewTrieRouter()
	err := router.RegisteredHandler(pattern)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	node := router.FindHandler(pattern)
	if node == nil {
		t.Fail()
		return
	}
	if node.isWild || node.pattern != pattern {
		t.Error("node pattern is: ", node.pattern)
		t.Fail()
	}
	t.Log(node)
}

func TestTrieRouter_RegisteredHandler_HalfPattern(t *testing.T) {
	pattern := "/hello/:name"
	router := NewTrieRouter()
	err := router.RegisteredHandler(pattern)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	url := "/hello/mike"
	node := router.FindHandler(url)
	if node == nil {
		t.Fail()
		return
	}
	if !node.isWild || node.part != ":name" {
		t.Error("node pattern is : ", node.pattern)
		t.Fail()
	}
	t.Log(node)
}

func TestTrieRouter_RegisteredHandler_AllPattern(t *testing.T) {
	patterns := []string{"/hello/mike", "/hello/:name/age", "/hello/*person"}
	router := NewTrieRouter()
	for _, pattern := range patterns {
		if err := router.RegisteredHandler(pattern); err != nil {
			t.Error(err)
			t.Fail()
			return
		}
	}

	url := "/hello/mike/male"
	node := router.FindHandler(url)
	if node == nil {
		t.Error("can not find a node")
		t.Fail()
		return
	}
	if !node.isWild || node.part != "*person" {
		t.Error("node pattern is : ", node.pattern)
		t.Fail()
	}
	t.Log(node)
}

func TestTrieRouter_RegisteredHandler_Priority(t *testing.T) {
	patterns := []string{"/hello/mike", "/hello/:name/age", "/hello/*person"}
	router := NewTrieRouter()
	for _, pattern := range patterns {
		if err := router.RegisteredHandler(pattern); err != nil {
			t.Error(err)
			t.Fail()
			return
		}
	}

	url1 := "/hello/mike"
	url2 := "/hello/mike/age"
	url3 := "/hello"

	if node := router.FindHandler(url1); node == nil || node.part != "mike" {
		t.Log(node)
		t.Fail()
	}
	if node := router.FindHandler(url2); node == nil || node.part != "age" {
		t.Log(node)
		t.Fail()
	}
	if node := router.FindHandler(url3); node != nil {
		t.Log(node)
		t.Fail()
	}
}