package cmd

import (
	"github.com/Prutswonder/submerger/merge"
	"github.com/spf13/cobra"
)

var (
	mergeCmd = &cobra.Command{
		Use:   "merge [path]",
		Short: "Merges all subs for the specified folder and its subfolders.",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runMerge(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(args []string) {
	path := ""

	if len(args) > 0 {
		path = args[0]
	}

	logger := merge.NewLogger()
	fileWalker := merge.NewFileWalker()
	commander := merge.NewCommander()
	merger := merge.NewMerger(supportedMovieExtensions,
		supportedSubtitleExtensions,
		mergedMovieExtension,
		logger,
		fileWalker,
		commander,
	)
	if err := merger.Run(path); err != nil {
		panic(err)
	}
}
