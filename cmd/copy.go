package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// copyCmd represents the vacuum command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy",
	Long:  `pg_stat_progress_copy.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.Copy)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
