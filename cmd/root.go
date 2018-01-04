package cmd

import (
	"fmt"

	"github.com/Prutswonder/submerger/merge"
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
			path := ""

			if len(args) > 0 {
				path = args[0]
			}

			merger := merge.NewMerger(supportedMovieExtensions, supportedSubtitleExtensions, mergedMovieExtension)
			if err := merger.Run(path); err != nil {
				panic(err)
			}
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
