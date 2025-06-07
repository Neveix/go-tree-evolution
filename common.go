package main

import (
	"fmt"
	"strings"
	"time"
)

func input(info string) string {
	fmt.Print(info + " > ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func ByteTo32mal(value byte) byte {
	if value <= 9 {
		return 48 + value
	} else if value <= 15 {
		return 65 + value - 10
	} else {
		return '_'
	}
}

func Now() int {
	now := time.Now()
	return int(now.UnixNano() / int64(time.Millisecond))
}

var sides = []int{1, 0, 0, 1, -1, 0, 0, -1}

func Remove[T comparable](slice []T, value T) []T {
	result := slice[:0]
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

func Repeat(n int, f func(int)) {
	for i := 0; i < n; i++ {
		f(i)
	}
}

func MergeLines(output1, output2, separator string) string {
	splitted1 := strings.Split(output1, "\n")
	splitted2 := strings.Split(output2, "\n")
	lineCount1 := len(splitted1)
	lineCount2 := len(splitted2)
	for i, line2 := range splitted2 {
		splitted1[lineCount1-lineCount2+i] += separator + line2
	}
	splitted1 = splitted1[:lineCount1-1]
	return strings.Join(splitted1, "\n")
}
