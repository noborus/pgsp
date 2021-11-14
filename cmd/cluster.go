package cmd

import (
	"github.com/noborus/pgsp/tui"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster",
	Long:  `pg_stat_progress_cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.Cluster)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
