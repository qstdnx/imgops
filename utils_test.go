package main

import (
	"testing"
)

func TestTargetSplit(t *testing.T) {
	queryList := getQueryList("google")
	if len(queryList) != 1 {
		t.Errorf("Expected 1 query: %v", queryList)
	}

	queryList = getQueryList("google, yandex, wrong")
	if len(queryList) != 2 {
		t.Error("Expected 2 queries")
	}
}

func TestNewTargets(t *testing.T) {
	queryList := getQueryList("iqdb, saucenao, tracemoe, ascii2d")
	if len(queryList) != 4 {
		t.Errorf("Expected 4 queries: %v", queryList)
	}

	m := getNameToIdTargets(availableTargets)
	expected := map[string]string{
		"iqdb":     "#t78",
		"saucenao": "#t82",
		"tracemoe": "#t201",
		"ascii2d":  "#t84",
	}

	for name, id := range expected {
		if got, ok := m[name]; !ok || got != id {
			t.Errorf("Expected '%s' to have ID '%s', got '%s'", name, id, got)
		}
	}
}

func TestAbsoluteURL(t *testing.T) {
	if got := absoluteURL("//imgops.com/a.png"); got != "https://imgops.com/a.png" {
		t.Errorf("Protocol-relative URL is wrong: %s", got)
	}

	if got := absoluteURL("/get2post?x=1"); got != "https://imgops.com/get2post?x=1" {
		t.Errorf("Relative URL is wrong: %s", got)
	}

	absolute := "https://iqdb.org/?url=x"
	if got := absoluteURL(absolute); got != absolute {
		t.Errorf("Absolute URL should be unchanged: %s", got)
	}
}

func TestGetUploadedImageURL(t *testing.T) {
	html := `<html><body><img id=mainImage rel=nofollow src="//imgops.com/1hr-tempcache/abc.png"></body></html>`

	url, err := getUploadedImageURL(html)
	if err != nil {
		t.Error(err)
	}

	if url != "https://imgops.com/1hr-tempcache/abc.png" {
		t.Errorf("Unexpected image URL: %s", url)
	}

	if _, err := getUploadedImageURL("<html></html>"); err == nil {
		t.Error("Expected an error when the image is missing")
	}
}
