package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

	arg "github.com/alexflint/go-arg"
)

var version = "1.0.3"

type args struct {
	Include []string `arg:"-i,separate,help:include lines which match patterns"`
	Exclude []string `arg:"-e,separate,help:exclude lines which match patterns"`
	Tail    []string `arg:"positional,help:highlight patters"`
}

func (args) Description() string {
	return "Highlights output from STDIN"
}

func (args) Version() string {
	return version
}

func main() {
	var args args
	arg.MustParse(&args)

	if len(args.Tail) > len(terminalColours) {
		fmt.Printf("Tail limit reached: %v \n", len(terminalColours))
		os.Exit(1)
	}

	tailRe := prepareRegExp(&args.Tail)
	includeRe := prepareRegExp(&args.Include)
	excludeRe := prepareRegExp(&args.Exclude)

	scanner := bufio.NewScanner(os.Stdin)

	process(scanner, &args.Tail, tailRe, includeRe, excludeRe)
}

// Process data from STDIN
func process(scanner *bufio.Scanner, argsTail *[]string, tailRe *regexp.Regexp, includeRe *regexp.Regexp, excludeRe *regexp.Regexp) {
	// Highlight matched pattern
	colourify := func(match string) string {
		index, err := getIndex(argsTail, match)
		if err != nil {
			return match
		}
		return terminalColours[index] + match + colourEnd
	}

	for scanner.Scan() {
		newLine := scanner.Text()

		// Check if the line should be included in the output
		if includeRe != nil {
			if !includeRe.MatchString(newLine) {
				continue
			}
		}

		// Check if the line should be excluded from the output
		if excludeRe != nil {
			if excludeRe.MatchString(newLine) {
				continue
			}
		}

		// Print original or colourified line
		if tailRe != nil {
			fmt.Print(tailRe.ReplaceAllStringFunc(newLine, colourify) + "\n")
		} else {
			fmt.Println(newLine)
		}
	}
}

// Returns nil or compiled regexp
func prepareRegExp(args *[]string) *regexp.Regexp {
	if len(*args) == 0 {
		return nil
	}
	return regexp.MustCompile(generateRegExp(args))
}

// Returns regular expression which is used for colourization
func generateRegExp(args *[]string) string {
	re := ""
	for _, value := range *args {
		if len(re) > 0 {
			re = re + "|"
		}
		re = re + "(" + regexp.QuoteMeta(value) + ")"
	}
	return re
}

// Returns position of the element in the array
func getIndex(array *[]string, element string) (int, error) {
	for index, value := range *array {
		if value == element {
			return index, nil
		}
	}
	return 0, errors.New(element + " not found")
}
