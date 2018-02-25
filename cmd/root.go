package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	supportedMovieExtensions    = []string{".avi", ".mkv", ".mp4"}
	supportedSubtitleExtensions = []string{".srt"}
	mergedMovieExtension        = "_subs.mkv"

	verbose = false

	rootCmd = &cobra.Command{
		Use:   "submerger [path]",
		Short: "Merge subtitles with movie files using MkvMerge.",
		Long: fmt.Sprintf("Merge subtitles with movie files using MkvMerge. "+
			"Supported movie extensions are: %v. "+
			"Supported subtitle extensions are: %v. "+
			"Merged file is written with extension [%v].",
			supportedMovieExtensions,
			supportedSubtitleExtensions,
			mergedMovieExtension,
		),
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runMerge(args)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
}

// Execute executes the root command
func Execute() {
	rootCmd.Execute()
}
