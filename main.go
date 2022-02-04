package main

import (
	"log"
	"os"
	"strings"

	"github.com/sgrech/parseargs"
)

const (
	GoogleBot = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	BingBot   = "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"
	YahooBot  = "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)"
)

func main() {
	c, err := parseargs.ParseArgs(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	baseUrl, ok := c.FindCommand("baseUrl")

	if !ok {
		log.Fatalln("--baseUrl is required")
	}

	entrypoint, _ := c.FindCommand("entrypoint")
	prefix, _ := c.FindCommand("prefix")
	userAgent, _ := c.FindCommand("userAgent")
	ignorePartial, ok := c.FindCommand("ignorePartial")

	var ignorePartialList []string

	if ok {
		ignorePartialList = strings.Split(ignorePartial, ",")
	}

	prefix = baseUrl + prefix
	uaMap := make(map[string]string)
	uaMap["google"] = GoogleBot
	uaMap["bing"] = BingBot
	uaMap["yahoo"] = YahooBot

	userAgent, ok = uaMap[userAgent]

	if !ok {
		userAgent = GoogleBot
	}

	us := NewUniqStack()
	us.Push(entrypoint)

	count := 0

	log.Printf("Starting fetch for %s with UserAgent %s", baseUrl+entrypoint, userAgent)
	for len(us.stack) > 0 {
		count++

		skipFilter := func(s string) bool {
			for _, v := range ignorePartialList {
				if strings.Contains(s, v) {
					return true
				}
			}
			return !strings.HasPrefix(s, prefix)
		}

		uri, err := us.Pop()
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		// if uri does not start with http, assume it is a relative
		// link and add the baseUrl to it
		if !strings.HasPrefix(uri, "http") {
			uri = baseUrl + uri
		}

		if skipFilter(uri) {
			log.Printf("%s %d %s %d %d/%d/%d", "SKIP", 0, uri, 0, count, len(us.hashSet), len(us.stack))
			continue
		}

		body, statusCode, err := GetContent(uri, userAgent)
		if err != nil {
			log.Printf("%s %d %s %d %d/%d/%d", "ERR", 0, uri, 0, count, len(us.hashSet), len(us.stack))
			continue
			// log.Fatalln(err)
			// os.Exit(1)
		}

		hrefs := GetHrefs(body)

		for _, href := range hrefs {
			us.Push(href)
		}
		log.Printf("%s %d %s %d %d/%d/%d", "GET", statusCode, uri, len(hrefs), count, len(us.hashSet), len(us.stack))
	}
}
