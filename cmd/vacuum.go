package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// vacuumCmd represents the vacuum command
var vacuumCmd = &cobra.Command{
	Use:   "vacuum",
	Short: "vacuum",
	Long:  `pg_stat_progress_vacuum.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.Vacuum)
	},
}

func init() {
	rootCmd.AddCommand(vacuumCmd)
}
