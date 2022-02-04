package main

import (
	"strings"
	"testing"
)

type getHrefsTest struct {
	body  string
	hrefs []string
	fn    hrefFilter
}

var htmlString = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
  <title>A Basic HTML5 Template</title>

  <link rel="icon" href="/favicon.ico">
  <link rel="icon" href="/favicon.svg" type="image/svg+xml">
  <link rel="apple-touch-icon" href="/apple-touch-icon.png">

  <link rel="stylesheet" href="css/styles.css?v=1.0">

</head>

<body>
  <!-- your content here... -->
	<a href="https://www.example.com/" target="_blank">Visit Example!</a> 
	<p><a href="https://www.example.org/">Example</a></p>
	<p><a href="https://www.example.com">Example 2</a></p>

	<h2>Relative URLs</h2>
	<p><a href="html_images.asp">HTML Images</a></p>
	<p><a href="/css/default.asp">CSS Tutorial</a></p>
  <script src="js/scripts.js"></script>
</body>
</html>
`

var getHrefsTests = []getHrefsTest{
	{
		"href=\"\"",
		[]string{},
		func(s string) bool {
			return true
		},
	},
	{
		"<a href=\"/en/test\"></a><a href=\"/test.js\"></a>",
		[]string{"/en/test"},
		func(s string) bool {
			return strings.HasPrefix(s, "/en")
		},
	},
	{
		"href=\"www.google.com\"",
		[]string{"www.google.com"},
		func(s string) bool {
			return true
		},
	},
	{
		"<a href=\"https://example.com\">Visit example.com!</a>",
		[]string{"https://example.com"},
		func(s string) bool {
			return true
		},
	},
	{
		htmlString,
		[]string{
			"/favicon.ico",
			"/favicon.svg",
			"/apple-touch-icon.png",
			"css/styles.css?v=1.0",
			"https://www.example.com/",
			"https://www.example.org/",
			"https://www.example.com",
			"html_images.asp",
			"/css/default.asp",
		},
		func(s string) bool {
			return true
		},
	},
}

func TestGetHrefs(t *testing.T) {
	for _, test := range getHrefsTests {
		hrefs := GetHrefs(test.body, test.fn)
		if len(hrefs) != len(test.hrefs) {
			t.Fatalf("Expected slice of length %d but got %d", len(test.hrefs), len(hrefs))
		}
		for i, href := range hrefs {
			if href != test.hrefs[i] {
				t.Fatalf("Expected value to be %s but got %s", test.hrefs[i], href)
			}
		}
	}
}
