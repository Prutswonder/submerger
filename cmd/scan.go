package cmd

import (
	"github.com/Prutswonder/submerger/merge"
	"github.com/spf13/cobra"
)

var (
	scanCmd = &cobra.Command{
		Use: "scan [path]",
		Short: "Returns a list of all merged movie files that also contain " +
			"the movie's original file.",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			path := ""

			if len(args) > 0 {
				path = args[0]
			}

			scanner := merge.NewScanner(supportedMovieExtensions, supportedSubtitleExtensions, mergedMovieExtension)
			if err := scanner.Run(path); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
}
