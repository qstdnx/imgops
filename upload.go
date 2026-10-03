package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const uploadUrl string = "https://imgops.com/store"
const uploadSearch string = "userUploadTempCache"
const defaultTarget = "imgops"

var defaultAction bool
var finalUrl string

func setDefaultAction(targetAction string) {
	if targetAction != defaultTarget {
		defaultAction = false
	} else {
		defaultAction = true
	}
}

func UploadURL(targetUrl string, targetAction string) ([]string, error) {
	setDefaultAction(targetAction)
	newUrl := "https://imgops.com/" + targetUrl
	finalUrl = newUrl

	if defaultAction {
		return strings.Fields(newUrl), nil
	} else {
		_, body, errs := gorequest.New().Get(newUrl).End()
		if errs != nil {
			return strings.Fields(""), errors.New("GET error")
		} else {
			return findHref(body, targetAction, finalUrl)
		}
	}
}

func redirectPolicy(req gorequest.Request, via []gorequest.Request) error {
	finalUrl = req.URL.String()
	if defaultAction {
		return errors.New("Stop redirection")
	} else {
		return nil
	}
}

func UploadFile(targetPath string, targetAction string) ([]string, error) {
	setDefaultAction(targetAction)

	debug("Read file from path")
	bytes, _ := ioutil.ReadFile(targetPath)

	debug("Make a POST request")
	_, body, errs := gorequest.New().
		Post(uploadUrl).
		RedirectPolicy(redirectPolicy).
		Type("multipart").
		SendFile(bytes, filepath.Base(targetPath), "photo").
		End()
	debug("POST request errors: %v", errs)

	debug("Final URL: %s", finalUrl)
	debug("Stop redirection: %b", defaultAction)
	if defaultAction {
		if len(finalUrl) > 0 {
			return strings.Fields(finalUrl), nil
		} else {
			return strings.Fields(""), errors.New("Could not get final URL after Redirect")
		}
	} else {
		return findHref(body, targetAction, finalUrl)
	}
}

// ImageURL returns the link of the image as stored by ImgOps.
// A local file is uploaded first, a remote URL is returned unchanged.
func ImageURL(targetPath string) (string, error) {
	if isUrl(targetPath) {
		return targetPath, nil
	}

	return UploadFileImageURL(targetPath)
}

// UploadFileImageURL uploads a local file to ImgOps and returns the URL of
// the uploaded image.
func UploadFileImageURL(targetPath string) (string, error) {
	debug("Read file from path")
	bytes, err := ioutil.ReadFile(targetPath)
	if err != nil {
		return "", err
	}

	debug("Make a POST request")
	_, body, errs := gorequest.New().
		Post(uploadUrl).
		Type("multipart").
		SendFile(bytes, filepath.Base(targetPath), "photo").
		End()
	debug("POST request errors: %v", errs)
	if len(errs) > 0 {
		return "", errors.New("POST error")
	}

	return getUploadedImageURL(body)
}
