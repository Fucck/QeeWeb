package router

import (
	"QeeWeb/qee/base"
	"net/http"
	"testing"
)

func TestTrieRouter_RegisteredHandler_absParse(t *testing.T) {
	pattern := "/hello/world"
	router := NewTrieRouter()
	f := func(ctx *base.Context){}
	err := router.RegisteredHandler(http.MethodGet, pattern, f)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	handler, _ := router.FindHandler(http.MethodGet, pattern)
	if handler == nil {
		t.Fail()
	}
}

func TestTrieRouter_RegisteredHandler_HalfPattern(t *testing.T) {
	pattern := "/hello/:name"
	router := NewTrieRouter()
	f := func(ctx *base.Context){}
	err := router.RegisteredHandler(http.MethodGet, pattern, f)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	url := "/hello/mike"
	handler, patternMap := router.FindHandler(http.MethodGet, url)
	if handler == nil {
		t.Fail()
		return
	}
	t.Log(patternMap)
	router.DisplayUrlPattern()
	if val, exist := patternMap["name"]; exist {
		if val != "mike" {
			t.Log("name: ", val, " error")
			t.Fail()
		}
	} else {
		t.Log("key name not exist")
	}
}

func TestTrieRouter_RegisteredHandler_AllPattern(t *testing.T) {
	patterns := []string{"/hello/mike", "/hello/:name/age", "/hello/*person"}
	router := NewTrieRouter()
	for _, pattern := range patterns {
		if err := router.RegisteredHandler(http.MethodGet, pattern, nil); err != nil {
			t.Error(err)
			t.Fail()
			return
		}
	}

	url := "/hello/mike/male"
	_, patternMap := router.FindHandler(http.MethodGet, url)
	t.Log("patternMap:", patternMap)
	if patternMap["person"] != "mike/male" {
		t.Log("person: ", patternMap["person"])
		t.Fail()
	}
}

func TestTrieRouter_RegisteredHandler_Priority(t *testing.T) {
	patterns := []string{"/hello/mike", "/hello/:name/age", "/hello/*person"}
	router := NewTrieRouter()
	for _, pattern := range patterns {
		if err := router.RegisteredHandler(http.MethodGet, pattern, nil); err != nil {
			t.Error(err)
			t.Fail()
			return
		}
	}

	url1 := "/hello/mike"
	url2 := "/hello/mike/age"
	url3 := "/hello/jack/class"

	if _, queryDict := router.FindHandler(http.MethodGet, url1); len(queryDict) != 0 {
		t.Log(queryDict)
		t.Fail()
	}
	if _, queryDict := router.FindHandler(http.MethodGet, url2); queryDict["name"] != "mike" {
		t.Log(queryDict)
		t.Fail()
	}
	if _, queryDict := router.FindHandler(http.MethodGet, url3); queryDict["person"] != "jack/class" {
		t.Log(queryDict)
		t.Fail()
	}
}