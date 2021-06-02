package cmd

import (
	"bufio"
	"fmt"
	"regexp"
)

// Returns list of regular expressions
func compileRegExp(args *[]string, regExpMode bool) ([]*regexp.Regexp, error) {
	all := make([]*regexp.Regexp, 0)
	for _, item := range *args {
		if !regExpMode {
			item = regexp.QuoteMeta(item)
		}
		pattern, err := regexp.Compile(item)
		if err != nil {
			return nil, err
		}
		all = append(all, pattern)
	}
	return all, nil
}

// Returns true if one of the regexp patterns matches the line
func matchesRegExpList(line *string, reList []*regexp.Regexp) bool {
	for _, item := range reList {
		if item.MatchString(*line) {
			return true
		}
	}
	return false
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
	includeRe []*regexp.Regexp, excludeRe []*regexp.Regexp, extract int, regexpMode bool) {

	for scanner.Scan() {
		newLine := scanner.Text()

		// Check if the line should be included in the output
		if len(includeRe) != 0 && !matchesRegExpList(&newLine, includeRe) {
			continue
		}

		// Check if the line should be excluded from the output
		if len(excludeRe) != 0 && matchesRegExpList(&newLine, excludeRe) {
			continue
		}

		// Apply limit
		if argsLimit != 0 {
			newLine = limitLine(&newLine, argsLimit)
		}

		// Print original or colorized line
		if colourifyRe != nil {
			if extract != 0 {
				for _, reItem := range colourifyRe {
					// Extract match
					result := reItem.FindAllString(newLine, -1)
					pointer := 0
					for _, i := range result {
						if i != "" {
							pointer++
							if pointer == extract {
								fmt.Println(i)
								break
							}
						}
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
