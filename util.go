package main

import (
	"io"
	"net/http"
	"regexp"
)

type hrefFilter func(string) bool

// Parse html body and extract all links enclosed in href quotes and match filter function fn
func GetHrefs(body string, fn hrefFilter) []string {
	r := regexp.MustCompile(`href="(.*?)"`)
	matches := r.FindAllStringSubmatch(body, -1)
	var hrefs []string
	for _, val := range matches {
		if len(val[1]) > 0 && fn(val[1]) { // check if href is not empty string and satisfies filter criteria
			hrefs = append(hrefs, val[1])
		}
	}
	return hrefs
}

// Fetch target url with google bot useragent
func GetContent(url string, userAgent string) (body string, statusCode int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	body = string(b)
	return body, resp.StatusCode, nil
}
