package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// basebackupCmd represents the basebackup command
var basebackupCmd = &cobra.Command{
	Use:   "basebackup",
	Short: "basebackup",
	Long:  `pg_stat_progress_basebackup.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.BaseBackup)
	},
}

func init() {
	rootCmd.AddCommand(basebackupCmd)
}
