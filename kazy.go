package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

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

func main() {
	args := os.Args[1:]
	re := regexp.MustCompile(generateRegExp(&args))
	scanner := bufio.NewScanner(os.Stdin)

	colourify := func(match string) string {
		index, err := getIndex(&args, match)
		if err != nil {
			return match
		}
		return terminalColours[index] + match + colourEnd
	}
	for scanner.Scan() {
		fmt.Printf(re.ReplaceAllStringFunc(scanner.Text(), colourify) + "\n")
	}
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
