package cmd

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
	info = "\nDescription:\n"
	info += fmt.Sprintf("     got length: %v\n", len(*result))
	info += fmt.Sprintf("expected length: %v\n", len(*expected))
	info += fmt.Sprintf("     got bytes: %v\n", *result)
	info += fmt.Sprintf("expected bytes: %v\n", *expected)
	info += fmt.Sprintf("     got strings: %s\n", *result)
	info += fmt.Sprintf("expected strings: %s\n", *expected)
	return info
}

func runProcess(
	scanner *bufio.Scanner,
	argsTail *[]string,
	argsLimit int,
	tailRe *regexp.Regexp,
	includeRe *regexp.Regexp,
	excludeRe *regexp.Regexp,
	extract bool,
) []byte {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	processData(scanner, argsTail, argsLimit, tailRe, includeRe, excludeRe, extract)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = old
	return out
}

func TestPassingString(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var argsLimit int

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestIncludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var argsLimit int
	includeRe := regexp.MustCompile("(1234)")

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, includeRe, nil, false)
	assertEqual(t, result, expected)
}

func TestExcludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var argsLimit int
	excludeRe := regexp.MustCompile("(1234)")

	expected := []byte("qwerty\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, nil, excludeRe, false)
	assertEqual(t, result, expected)
}

func TestColourifyString(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{input}
	var argsLimit int
	tailRe := regexp.MustCompile("(1234)")

	expected := []byte("\033[46m1234\033[0m\n")

	result := runProcess(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestColourifyMultiple1(t *testing.T) {
	const input = "Jun 05 18:17:32 dell firefox.desktop[4089]: onEvent@resource://gre/modules/commonjs/toolkit/loader.js"
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{"5", "dell", "firefox", "403", "modules/", "loader.js"}
	var argsLimit int
	tailRe := prepareRegExp(&argsTail)

	expected := []byte("Jun 0\033[46m5\033[0m 18:17:32 \033[41mdell\033[0m \033[42mfirefox\033[0m.desktop[4089]: onEvent@resource://gre/\033[44mmodules/\033[0mcommonjs/toolkit/\033[45mloader.js\033[0m\n")

	result := runProcess(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestColourifyMultiple2(t *testing.T) {
	const input = "1 2 1 2"
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{"2"}
	var argsLimit int
	tailRe := prepareRegExp(&argsTail)

	expected := []byte("1 \033[46m2\033[0m 1 \033[46m2\033[0m\n")

	result := runProcess(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestColourifyPercentString(t *testing.T) {
	const input = "%"
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{input}
	var tailRe *regexp.Regexp
	var argsLimit int
	tailRe = regexp.MustCompile("(%)")

	expected := []byte("\033[46m%\033[0m\n")

	result := runProcess(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestColourifySquareBracketString(t *testing.T) {
	const input = "["
	scanner := bufio.NewScanner(strings.NewReader(input))
	argsTail := []string{input}
	var tailRe *regexp.Regexp
	var argsLimit int
	tailRe = prepareRegExp(&argsTail)

	expected := []byte("\033[46m[\033[0m\n")

	result := runProcess(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestExcludeIncludeString(t *testing.T) {
	const input = "1234\nqwerty"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	var argsLimit int
	var excludeRe *regexp.Regexp
	var includeRe *regexp.Regexp
	excludeRe = regexp.MustCompile("(1234)")
	includeRe = regexp.MustCompile("(1234)")

	expected := []byte("")

	result := runProcess(scanner, &argsTail, argsLimit, nil, includeRe, excludeRe, false)
	assertEqual(t, result, expected)
}

func TestEscapeInRegexp(t *testing.T) {
	test := []string{"["}
	prepareRegExp(&test)
}

func TestLimitStringSmaller(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	argsLimit := 2

	expected := []byte("12...\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestLimitStringEqual(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	argsLimit := 4

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestLimitStringBigger(t *testing.T) {
	const input = "1234"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var argsTail []string
	argsLimit := 10

	expected := []byte("1234\n")

	result := runProcess(scanner, &argsTail, argsLimit, nil, nil, nil, false)
	assertEqual(t, result, expected)
}

func TestExtractSimplePresent(t *testing.T) {
	const input = "    ↳ Microsoft Microsoft® Nano Transceiver v2.0	id=10	[slave  keyboard (3)]"
	scanner := bufio.NewScanner(strings.NewReader(input))
	tail := []string{"3"}
	tailRe := prepareRegExp(&tail)
	result := runProcess(scanner, &tail, 0, tailRe, nil, nil, true)
	assertEqual(t, result, []byte("3\n"))
}

func TestExtractSimpleWord(t *testing.T) {
	const input = "    ↳ Microsoft Microsoft® Nano Transceiver v2.0	id=10	[slave  keyboard (3)]"
	scanner := bufio.NewScanner(strings.NewReader(input))
	tail := []string{"id"}
	tailRe := prepareRegExp(&tail)
	result := runProcess(scanner, &tail, 0, tailRe, nil, nil, true)
	assertEqual(t, result, []byte("id\n"))
}

func BenchmarkProcess(b *testing.B) {
	const sample = "Jun 05 18:17:32 dell firefox.desktop[4089]: onEvent@resource://gre/modules/commonjs/toolkit/loader.js"

	scanner := bufio.NewScanner(strings.NewReader(sample))
	var argsTail []string
	var argsLimit int
	tailRe := regexp.MustCompile("(5)|(firefox)|(dell)|(o)|(s)")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processData(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	}
}

func BenchmarkProcessWithLimit(b *testing.B) {
	const sample = "Jun 05 18:17:32 dell firefox.desktop[4089]: onEvent@resource://gre/modules/commonjs/toolkit/loader.js"

	scanner := bufio.NewScanner(strings.NewReader(sample))
	var argsTail []string
	argsLimit := 50
	tailRe := regexp.MustCompile("(5)|(firefox)|(dell)|(o)|(s)")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processData(scanner, &argsTail, argsLimit, tailRe, nil, nil, false)
	}
}
