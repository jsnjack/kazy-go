package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

	arg "github.com/alexflint/go-arg"
)

var version = "1.0.1"
var colourEnd = "\033[0m"

var terminalColours = []string{
	"\033[46m",
	"\033[41m",
	"\033[42m",
	"\033[43m",
	"\033[44m",
	"\033[45m",
	"\033[101m",
	"\033[102m",
	"\033[104m",
	"\033[105m",
	"\033[106m",
	"\033[31m",
	"\033[32m",
	"\033[33m",
	"\033[34m",
	"\033[35m",
	"\033[36m",
	"\033[91m",
	"\033[92m",
	"\033[93m",
	"\033[94m",
	"\033[95m",
	"\033[96m",
}

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
			fmt.Printf(tailRe.ReplaceAllStringFunc(newLine, colourify) + "\n")
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
		re = re + "(" + value + ")"
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
