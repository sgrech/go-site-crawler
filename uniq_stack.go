package main

import "errors"

type UniqStack struct {
	stack   []string
	hashSet map[string]int
}

func NewUniqStack() *UniqStack {
	s := []string{}
	hs := make(map[string]int)
	return &UniqStack{s, hs}
}

func (us *UniqStack) Push(val string) {
	valInHashSet := us.hashSet[val]
	if valInHashSet == 0 {
		us.stack = append(us.stack, val)
		us.hashSet[val] = 1
	}
}

func (us *UniqStack) Pop() (val string, err error) {
	if len(us.stack) == 0 {
		return "", errors.New("Cannot pop empty stack")
	}
	n := len(us.stack) - 1
	val = us.stack[n]
	us.stack = us.stack[:n]
	return val, nil
}
