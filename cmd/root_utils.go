package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
)

// Returns nil or compiled regexp
func prepareRegExp(args *[]string) *regexp.Regexp {
	if len(*args) == 0 {
		return nil
	}
	return regexp.MustCompile(generateRegExp(args))
}

// Returns regular expression which is used for colorization
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

// limitLine limits length of the line
func limitLine(line *string, limit int) string {
	const marker = "..."
	if len(*line) > limit {
		var l string
		l = *line
		return l[:limit] + "..."
	}
	return *line
}

// Process data from STDIN
func processData(scanner *bufio.Scanner, argsTail *[]string, argsLimit int, tailRe *regexp.Regexp,
	includeRe *regexp.Regexp, excludeRe *regexp.Regexp) {
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

		// Apply limit
		if argsLimit != 0 {
			newLine = limitLine(&newLine, argsLimit)
		}

		// Print original or colorized line
		if tailRe != nil {
			fmt.Print(tailRe.ReplaceAllStringFunc(newLine, colourify) + "\n")
		} else {
			fmt.Println(newLine)
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}
