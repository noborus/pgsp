package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/noborus/pgsp"
	"github.com/noborus/pgsp/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	afterCompletion int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgsp",
	Short: "pg_stat_progress monitor",
	Long: `Monitors PostgreSQL's pg_stat_progress_*.
`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := pgsp.Connect()
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()

		model := tui.Model{
			DB: db,
		}
		tui.AfterCompletion = time.Duration(afterCompletion)

		p := tea.NewProgram(model)
		if err := p.Start(); err != nil {
			fmt.Printf("there's been an error: %v", err)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pspt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().IntVarP(&afterCompletion, "AfterCompletion", "a", 10, "Number of seconds to display after completion")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pspt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pspt")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
