package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// createindexCmd represents the createindex command
var createindexCmd = &cobra.Command{
	Use:   "createindex",
	Short: "createindex",
	Long:  `pg_stat_progress_createindex.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.CreateIndex)
	},
}

func init() {
	rootCmd.AddCommand(createindexCmd)
}
