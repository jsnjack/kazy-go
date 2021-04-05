package cmd

import (
	"bufio"
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

// Returns list of regular expressions
func compileRegExp(args *[]string) []*regexp.Regexp {
	all := make([]*regexp.Regexp, 0)
	for _, item := range *args {
		all = append(all, regexp.MustCompile(regexp.QuoteMeta(item)))
	}
	return all
}

// limitLine limits length of the line
func limitLine(line *string, limit int) string {
	if len(*line) > limit {
		l := *line
		return l[:limit] + "..."
	}
	return *line
}

// Process data from STDIN
func processData(scanner *bufio.Scanner, argsLimit int, colourifyRe []*regexp.Regexp,
	includeRe *regexp.Regexp, excludeRe *regexp.Regexp, extract bool) {

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
		if colourifyRe != nil {
			if extract {
				var match string
				for _, reItem := range colourifyRe {
					match = reItem.FindString(newLine)
					if match != "" {
						fmt.Println(match)
						break
					}
				}
			} else {
				for idx, reItem := range colourifyRe {
					colourify := func(match string) string {
						return terminalColours[idx] + match + colourEnd
					}
					newLine = reItem.ReplaceAllStringFunc(newLine, colourify)
				}
				fmt.Println(newLine)
			}
		} else {
			fmt.Println(newLine)
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}
