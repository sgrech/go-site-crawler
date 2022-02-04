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
	baseUrl, found := c.FindCommand("baseUrl")

	if !found {
		log.Fatalln("--baseUrl is required")
	}

	entrypoint, _ := c.FindCommand("entrypoint")
	prefix, _ := c.FindCommand("prefix")
	userAgent, _ := c.FindCommand("userAgent")

	uaMap := make(map[string]string)
	uaMap["google"] = GoogleBot
	uaMap["bing"] = BingBot
	uaMap["yahoo"] = YahooBot

	userAgent, ok := uaMap[userAgent]

	if !ok {
		userAgent = GoogleBot
	}

	us := NewUniqStack()
	us.Push(entrypoint)

	count := 0

	log.Printf("Starting fetch for %s with UserAgent %s", baseUrl+entrypoint, userAgent)
	for len(us.stack) > 0 {
		count++
		uri, err := us.Pop()
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		body, statusCode, err := GetContent(baseUrl+uri, userAgent)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		filter := func(s string) bool {
			return strings.HasPrefix(s, prefix)
		}

		hrefs := GetHrefs(body, filter)

		for _, href := range hrefs {
			us.Push(href)
		}
		log.Printf("%d %s %d %d/%d - %d", statusCode, baseUrl+uri, len(hrefs), count, len(us.hashSet), len(us.stack))
	}
}
