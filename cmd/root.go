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

// Version is the version of the application calculated with monova
var Version string

var rootCmd = &cobra.Command{
	Use:   "kazy [<pattern>...]",
	Short: "Highlights, filters and extracts string patterns from STDIN",
	Args: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		if len(args) > len(terminalColours) {
			return fmt.Errorf("tail limit reached: %d", len(terminalColours))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if rootVersion {
			fmt.Println(Version)
			return nil
		}

		tailRe := prepareRegExp(&args)
		includeRe := prepareRegExp(&rootInclude)
		excludeRe := prepareRegExp(&rootExclude)

		scanner := bufio.NewScanner(os.Stdin)
		// Update max string size from 64 to 1024
		buffer := make([]byte, 0, 64*1024)
		scanner.Buffer(buffer, 1024*1024)

		processData(scanner, &args, rootLimit, tailRe, includeRe, excludeRe, rootExtractMode)

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
	rootCmd.PersistentFlags().StringArrayVarP(
		&rootInclude, "include", "i", rootInclude, "only include lines which match provided patterns",
	)
	rootCmd.PersistentFlags().StringArrayVarP(
		&rootExclude, "exclude", "e", rootExclude, "exclude from output lines which match provided patterns",
	)
	rootCmd.PersistentFlags().BoolVar(&rootVersion, "version", false, "print version and exit")
	rootCmd.PersistentFlags().BoolVarP(&rootExtractMode, "extract", "x", false, "extract matched strings instead of highlighting them")
}
