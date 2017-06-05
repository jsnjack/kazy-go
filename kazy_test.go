package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, result []byte, expected []byte) {
	if len(result) == len(expected) {
		for index, value := range result {
			if value != expected[index] {
				t.Error(createInfo(&result, &expected))
				break
			}
		}
	} else {
		t.Errorf(createInfo(&result, &expected))
	}
}

func createInfo(result *[]byte, expected *[]byte) string {
	var info string
	info = "Result vs expected:\n"
	info += fmt.Sprintf("length: %v vs %v\n", len(*result), len(*expected))
	info += fmt.Sprintf("bytes: %v vs %v\n", *result, *expected)
	info += fmt.Sprintf("bytes: %s vs %s\n", *result, *expected)
	return info
}

func runProcess(
	scanner *bufio.Scanner,
	argsTail *[]string,
	tailRe *regexp.Regexp,
	includeRe *regexp.Regexp,
	excludeRe *regexp.Regexp,
) []byte {
	r, w, _ := os.Pipe()
	os.Stdout = w

	process(scanner, argsTail, tailRe, includeRe, excludeRe)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	return out
}

func TestPassingString(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, nil, nil, nil)
	assertEqual(t, result, expected)
}

func TestIncludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var includeRe *regexp.Regexp
	includeRe = regexp.MustCompile("(1234)")

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, nil, includeRe, nil)
	assertEqual(t, result, expected)
}

func TestExcludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var excludeRe *regexp.Regexp
	excludeRe = regexp.MustCompile("(1234)")

	expected := []byte("qwerty\n")

	result := runProcess(scanner, &argsTail, nil, nil, excludeRe)
	assertEqual(t, result, expected)
}

func TestColourifyString(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{input}
	var tailRe *regexp.Regexp
	tailRe = regexp.MustCompile("(1234)")

	expected := []byte("\033[46m1234\033[0m\n")

	result := runProcess(scanner, &argsTail, tailRe, nil, nil)
	assertEqual(t, result, expected)
}

func TestExcludeIncludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var excludeRe *regexp.Regexp
	var includeRe *regexp.Regexp
	excludeRe = regexp.MustCompile("(1234)")
	includeRe = regexp.MustCompile("(1234)")

	expected := []byte("")

	result := runProcess(scanner, &argsTail, nil, includeRe, excludeRe)
	assertEqual(t, result, expected)
}

func BenchmarkProcess(b *testing.B) {
	const sample = "Jun 05 18:17:32 dell firefox.desktop[4089]: onEvent@resource://gre/modules/commonjs/toolkit/loader.js"

	scanner := bufio.NewScanner(strings.NewReader(sample))
	var argsTail []string
	var tailRe *regexp.Regexp
	tailRe = regexp.MustCompile("(5)|(firefox)|(dell)|(o)|(s)")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		process(scanner, &argsTail, tailRe, nil, nil)
	}
}
