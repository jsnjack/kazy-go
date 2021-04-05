package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootLimit int
var rootInclude []string
var rootExclude []string
var rootVersion bool
var rootExtractMode bool
var rootRegExpMode bool
var bufferSize int

// Version is the version of the application calculated with monova
var Version string

var rootCmd = &cobra.Command{
	Use:   "kazy [<pattern>...]",
	Short: "Highlights, filters and extracts string patterns from STDIN",
	Args: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		if len(args) > len(terminalColours) {
			return fmt.Errorf("pattern limit reached: %d", len(terminalColours))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if rootVersion {
			fmt.Println(Version)
			return nil
		}

		colourifyRe, err := compileRegExp(&args, rootRegExpMode)
		if err != nil {
			return err
		}

		includeRe, err := compileRegExp(&rootInclude, rootRegExpMode)
		if err != nil {
			return err
		}

		excludeRe, err := compileRegExp(&rootExclude, rootRegExpMode)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)
		buffer := make([]byte, 0, bufferSize*1024)
		scanner.Buffer(buffer, bufferSize*1024)

		processData(scanner, rootLimit, colourifyRe, includeRe, excludeRe, rootExtractMode, rootRegExpMode)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&rootLimit, "limit", "l", rootLimit, "limit the length of the line, characters")
	rootCmd.PersistentFlags().IntVarP(&bufferSize, "buffer", "b", 64, "buffer size in KB")
	rootCmd.PersistentFlags().StringArrayVarP(
		&rootInclude, "include", "i", rootInclude, "only include lines which match provided patterns",
	)
	rootCmd.PersistentFlags().StringArrayVarP(
		&rootExclude, "exclude", "e", rootExclude, "exclude from output lines which match provided patterns",
	)
	rootCmd.PersistentFlags().BoolVar(&rootVersion, "version", false, "print version and exit")
	rootCmd.PersistentFlags().BoolVarP(
		&rootExtractMode, "extract", "x", false,
		"extract matched strings (leftmost) instead of highlighting them",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&rootRegExpMode, "regexp", "r", false,
		"use RegExp patterns instead of string patterns",
	)
}
