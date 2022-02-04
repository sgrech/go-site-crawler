# Go Site Crawler

Go Site Crawler is a simple application written in go that can fetch content
from a url endpoint, scan the content for href links and subsequently crawl the
entire site. This is not done for data collection or scraping purposes and it
does little more then fetch content and crawl hrefs.

It has a log output that will tell you the data and time the content has been
fetched, the status code it received from the endpoint it tried to query, the
number of hrefs it found and how many links it has parsed so far out of how many
remain to be parsed.

The application uses a stack that tracks which endpoints have already been added
so as to avoid adding duplicates. This should ensure that, even if the same href
is present multiple times in the content, it will only fetch that endpoint once.

There are a number of reasons why you might want to do this:

1. You want to traverse the entire site and ensure that there are no broken
   links (will give a 4xx or 5xx status code response)
2. You want to traverse the entire site so as to build cache (such as for
   example prerender.io or you have a site that has some kind of caching
   mechanism)
3. You want to generate a list of all the links that can be detected on your
   website.

## Building the application

You will need go installed on your target os to build this application, the
current release has been built using go version 1.17.5 but it will probably work
for older versions of go too.

You can build the application by cloning this repo and in it's project folder
running the following:

```
go build .
```

The process should generate an executable called `go-site-crawler`.

## Running the application

You can run the application after building it by using the following command:

```
./go-site-crawler --baseUrl=https://example.com
```

The full list of arguments is as follows

```
--baseUrl=https://exmaple.com - this is the base endpoint that will be used when
crawling the site

--entrypoint=/ - this is the entrypint that is used in conjunction with the
baseUrl to fetch content. If omitted it will resolve to an empty string and
first endpoint to be fetched will be the baseUrl

--prefix=/en - This is a filter that will be used when scanning for href
content. Giving for example a prefix of `/` will only scan for relative urls,
while giving a prefix of `/en` will only scan for urls that start with `/en`

--userAgent=google - This is the user agent you want the crawler to use when
fetching content. Currently the only ones available are 'google', 'bing' and
'yahoo', and it will use the respective bot user agent. This defautls to google
bot.
```

## TODO

Will be adding a Dockerfile so that the application can be run via docker so as
to avoid having to install go if you really don't need it.
