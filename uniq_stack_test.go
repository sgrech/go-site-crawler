package main

import (
	"testing"
)

type uniqStackPushTest struct {
	args               []string
	expectedStack      []string
	expectedHashSetLen int
}

var uniqStackPushTests = []uniqStackPushTest{
	{
		[]string{"test"},
		[]string{"test"},
		1,
	},
	{
		[]string{"test", "test2", "test3"},
		[]string{"test", "test2", "test3"},
		3,
	},
	{
		[]string{"test", "test2", "test"},
		[]string{"test", "test2"},
		2,
	},
	{
		[]string{"test", "test2", "test2", "test"},
		[]string{"test", "test2"},
		2,
	},
}

func TestUniqStackPush(t *testing.T) {
	for _, test := range uniqStackPushTests {
		us := NewUniqStack()
		for _, arg := range test.args {
			us.Push(arg)
		}
		if len(us.stack) != len(test.expectedStack) {
			t.Fatalf("Expected stack of size %d but got %d", len(test.expectedStack), len(us.stack))
		}

		if len(us.hashSet) != test.expectedHashSetLen {
			t.Fatalf("Expected hashSet of size %d but got %d", test.expectedHashSetLen, len(us.hashSet))
		}

		for i, val := range us.stack {
			if val != test.expectedStack[i] {
				t.Fatalf("Expect element %d to be %s but got %s", i, test.expectedStack[i], val)
			}
		}
	}
}

type uniqStackPopTest struct {
	args               []string
	numOfPops          int
	expectedVals       []string
	expectedStackLen   int
	expectedHashSetLen int
	expectedError      bool
}

var uniqStackPopTests = []uniqStackPopTest{
	{
		[]string{},
		1,
		[]string{},
		0,
		0,
		true,
	},
	{
		[]string{"test"},
		1,
		[]string{"test"},
		0,
		1,
		false,
	},
	{
		[]string{"test", "test2"},
		1,
		[]string{"test2"},
		1,
		2,
		false,
	},
	{
		[]string{"test", "test2", "test3", "test2"},
		1,
		[]string{"test3"},
		2,
		3,
		false,
	},
}

func TestUniqStackPop(t *testing.T) {
	for _, test := range uniqStackPopTests {
		us := NewUniqStack()
		for _, arg := range test.args {
			us.Push(arg)
		}

		for i := 0; i < test.numOfPops; i++ {
			val, err := us.Pop()
			if err != nil && !test.expectedError {
				t.Fatalf("Pop generated an error when it was not expected to")
			} else if err != nil && test.expectedError {
				continue
			}
			if val != test.expectedVals[i] {
				t.Fatalf("Expect element %d to be %s but got %s", i, test.expectedVals[i], val)
			}
		}
		if len(us.stack) != test.expectedStackLen {
			t.Fatalf("Expected stack of size %d but got %d", test.expectedStackLen, len(us.stack))
		}

		if len(us.hashSet) != test.expectedHashSetLen {
			t.Fatalf("Expected hashSet of size %d but got %d", test.expectedHashSetLen, len(us.hashSet))
		}

	}
}
