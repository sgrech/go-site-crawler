package main

import (
	"errors"
	"net/url"
	"strings"
)

type UrlFilter struct {
	baseUrl    string
	rootUrl    string
	conditions []Condition
}

// type Conditioner interface {
// 	Condition(testUrl string) bool
// }
type Condition func(testUrl string) bool

func NewUrlFilter(baseUrl string) (uf *UrlFilter, err error) {
	u, err := url.Parse(baseUrl)

	emptyTests := []Condition{}

	if err != nil {
		return &UrlFilter{"", "", emptyTests}, err
	}

	if u.Host == "" {
		return &UrlFilter{"", "", emptyTests}, errors.New("Parsing url yielded no host, assuming url is invalid")
	}

	rootUrl := u.Scheme + "://" + u.Host

	return &UrlFilter{baseUrl, rootUrl, emptyTests}, nil
}

func (uf UrlFilter) IsValid(url string) (ok bool, testUrl string) {
	// if url is a relative one (starts with a forward slash), then prefix it with the base url
	if strings.HasPrefix(url, "/") {
		testUrl = uf.baseUrl + url
	} else {
		testUrl = url
	}

	for _, condition := range uf.conditions {
		if ok := !condition(testUrl); !ok {
			return ok, testUrl
		}
	}

	return true, testUrl
}

func (uf *UrlFilter) AppendCondition(c Condition) {
	uf.conditions = append(uf.conditions, c)
}

func (uf *UrlFilter) PrefixCondition() Condition {
	return func(testUrl string) bool {
		return strings.HasPrefix(testUrl, uf.baseUrl)
	}
}
