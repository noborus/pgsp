package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "analyze",
	Long:  `pg_stat_progress_analyze.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.Analyze)
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
