package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isUrl(targetPath string) bool {
	return strings.Contains(targetPath, "http:") || strings.Contains(targetPath, "https:")
}

// absoluteURL turns links returned by ImgOps into absolute URLs.
// ImgOps links are usually absolute but some targets (e.g. ascii2d) use
// site-relative links such as "/get2post?...".
func absoluteURL(href string) string {
	href = strings.TrimSpace(href)

	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	if strings.HasPrefix(href, "/") {
		return "https://imgops.com" + href
	}

	return href
}

// getUploadedImageURL extracts the URL of the image shown by ImgOps from a
// result page. For uploaded files this is the temporary ImgOps cache URL,
// for remote images it is the original URL.
func getUploadedImageURL(document string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return "", err
	}

	src, attrOk := doc.Find("#mainImage").Attr("src")
	if !attrOk || strings.TrimSpace(src) == "" {
		return "", errors.New("Could not find the image URL in the ImgOps response")
	}

	return absoluteURL(src), nil
}

func findHref(document, targetStr, finalUrl string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return strings.Fields(""), err
	} else {
		queryList := getQueryList(targetStr)
		debug("Queries: %v", queryList)
		foundUrls := make([]string, 0, len(queryList))

		for _, query := range queryList {
			href, attrOk := doc.Find(query).Attr("href")
			debug("Get href from query '%s': %v", query, attrOk)

			if attrOk {
				foundUrls = append(foundUrls, absoluteURL(href))
			}
		}

		if len(foundUrls) > 0 {
			return foundUrls, nil
		} else {
			return strings.Fields(finalUrl), errors.New(fmt.Sprintf("No link is found in targets: %s", queryList))
		}
	}
}

// Returns queries of found targets
func getQueryList(s string) []string {
	split := strings.Split(s, ",")
	m := getNameToIdTargets(availableTargets)

	found := make([]string, 0, len(split))
	for _, val := range split {
		key := strings.TrimSpace(val)
		query, mapOk := m[key]
		if mapOk {
			debug("Query for '%s' is found: %s", key, query)
			found = append(found, query)
		}
	}

	return found
}

// Create a nice select list

func genSelectText(inputMode bool) string {
	text := ""

	if inputMode {
		text += "Available targets:\n\n"
	} else {
		text += "Select a target:\n\n"
	}

	text += "  > [i]mgops\n"
	for _, target := range availableTargets {
		text += "  > " + strings.Replace(target.Name, string(target.Key), "["+string(target.Key)+"]", 1) + "\n"
	}

	if inputMode {
		text += "\nType initials: "
	} else {
		text += "\n(Press ESC to cancel)"
	}

	return text
}

// Initials to targets with comma

func initialsToTargets(initials string) string {
	targets := ""
	m := getKeyToNameTargets(availableTargets)

	for i := range initials {
		if val, ok := m[rune(initials[i])]; ok {
			targets += val
			if len(initials)-1 != i {
				targets += ","
			}
		}
	}

	return targets
}
