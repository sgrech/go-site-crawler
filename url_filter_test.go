package main

import (
	"strings"
	"testing"
)

type newUrlFilterTest struct {
	url             string
	expectedRootUrl string
	expectedErr     bool
}

var newUrlFilterTests = []newUrlFilterTest{
	{
		"https://example.com",
		"https://example.com",
		false,
	},
	{
		"https://example.com/en",
		"https://example.com",
		false,
	},
	{
		"https://www.example.com/en",
		"https://www.example.com",
		false,
	},
	{
		"http://example.com",
		"http://example.com",
		false,
	},
	{
		"http://example.com/en",
		"http://example.com",
		false,
	},
	{
		"htt/example.com/en",
		"",
		true,
	},
	{
		"example.com/en",
		"",
		true,
	},
}

func TestNewUrlFilter(t *testing.T) {
	for _, test := range newUrlFilterTests {
		uf, err := NewUrlFilter(test.url)

		if err != nil && !test.expectedErr {
			t.Fatalf("Expected no err but got \"%s\"", err)
		} else if err == nil && test.expectedErr {
			t.Fatal("Expected err but got none")
		} else if err != nil && test.expectedErr {
			// If testing for an error then skip the next tests
			// as they will fail which is expected but should
			// not fail the test
			continue
		}

		if uf.rootUrl != test.expectedRootUrl {
			t.Fatalf("Expected root url %s but got %s", test.expectedRootUrl, uf.rootUrl)
		}

		if uf.baseUrl != test.url {
			t.Fatalf("Expected base url %s but got %s", test.url, uf.baseUrl)
		}
	}
}

type isValidTest struct {
	baseUrl         string
	url             string
	expectedOk      bool
	expectedTestUrl string
}

var isValidTests = []isValidTest{
	{
		"https://example.com",
		"https://example.com",
		true,
		"https://example.com",
	},
	{
		"https://example.com",
		"/en",
		true,
		"https://example.com/en",
	},
}

func TestIsValid(t *testing.T) {
	for _, test := range isValidTests {
		uf, err := NewUrlFilter(test.baseUrl)

		if err != nil {
			t.Fatal(err)
		}

		ok, tr := uf.IsValid(test.url)

		if !strings.HasPrefix(tr, test.baseUrl) {
			t.Fatalf("Expected testUrl to have %s prefix but got %s", test.baseUrl, tr)
		}

		if ok != test.expectedOk {
			t.Fatalf("Expected valid result to be %v but got %v", test.expectedOk, ok)
		}
	}
}

type appendConditionTest struct {
	baseUrl                  string
	conditions               []Condition
	expectedConditionsLen    int
	expectedConditionResults []bool
}

var appendConditionTests = []appendConditionTest{
	{
		"https://example.com",
		[]Condition{
			func(testUrl string) bool { return true },
		},
		1,
		[]bool{true},
	},
	{
		"https://example.com",
		[]Condition{
			func(testUrl string) bool { return true },
			func(testUrl string) bool { return false },
		},
		2,
		[]bool{true, false},
	},
	{
		"https://example.com",
		[]Condition{
			func(testUrl string) bool { return false },
			func(testUrl string) bool { return true },
			func(testUrl string) bool { return true },
		},
		3,
		[]bool{false, true, true},
	},
	{
		"https://example.com",
		[]Condition{
			func(testUrl string) bool { return true },
			func(testUrl string) bool { return false },
			func(testUrl string) bool { return true },
			func(testUrl string) bool { return true },
			func(testUrl string) bool { return false },
		},
		5,
		[]bool{true, false, true, true, false},
	},
}

func TestAppendCondition(t *testing.T) {
	for _, test := range appendConditionTests {
		uf, err := NewUrlFilter(test.baseUrl)

		if err != nil {
			t.Fatal(err)
		}

		for _, condition := range test.conditions {
			uf.AppendCondition(condition)
		}

		if len(uf.conditions) != test.expectedConditionsLen {
			t.Fatalf("Expected conditions slice of len %d but got %d", test.expectedConditionsLen, len(uf.conditions))
		}

		if len(uf.conditions) != len(test.expectedConditionResults) {
			t.Fatalf("Number of conditions appended %d does not match number of expected results %d", len(uf.conditions), len(test.expectedConditionResults))
		}

		for i, condition := range uf.conditions {
			if condition("") != test.expectedConditionResults[i] {
				t.Fatalf("Expected condition result %v but instead got %v", test.expectedConditionResults[i], condition(""))
			}
		}
	}
}

type prefixConditionTest struct {
	baseUrl        string
	url            string
	expectedResult bool
}

var prefixConditionTests = []prefixConditionTest{
	{
		"https://example.com",
		"https://example.com",
		true,
	},
	{
		"https://example.com",
		"https://example.com/en",
		true,
	},
	{
		"https://example.com/en",
		"https://example.com/en",
		true,
	},
	{
		"https://example.com/en",
		"https://example.com/en/test/help",
		true,
	},
	{
		"http://example.com",
		"https://example.com/en",
		false,
	},
	{
		"https://example.com/fi",
		"https://ex.com/en",
		false,
	},
	{
		"https://example.com/fi",
		"https://example.com/en",
		false,
	},
}

func TestPrefixCondition(t *testing.T) {
	for _, test := range prefixConditionTests {
		uf, err := NewUrlFilter(test.baseUrl)

		if err != nil {
			t.Fatal(err)
		}

		pc := uf.PrefixCondition()
		result := pc(test.url)

		if result != test.expectedResult {
			t.Fatalf("Expected condition result to be %v but got %v", test.expectedResult, result)
		}
	}
}
