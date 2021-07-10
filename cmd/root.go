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
	cfgFile string
	config  Config
)

type Config struct {
	DSN             string  `yaml:"dsn"`
	AfterCompletion int     `yaml:"AfterCompletion"`
	Interval        float64 `yaml:"Interval"`
}

var (
	version         bool
	dsn             string
	afterCompletion int
	interval        float64
)

var (
	// Version represents the version.
	Version = "dev"
	// Revision set "git rev-parse --short HEAD".
	Revision = "HEAD"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgsp",
	Short: "pg_stat_progress monitor",
	Long: `Monitors PostgreSQL's pg_stat_progress_*.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("pgsp version %s rev:%s\n", Version, Revision)
			return
		}

		setConfig(config)
		tui.AfterCompletion = time.Duration(afterCompletion)
		tui.UpdateInterval = time.Duration(time.Second * time.Duration(interval))

		db, err := pgsp.Connect(dsn)
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()

		model := tui.NewModel(db)

		p := tea.NewProgram(model)
		if err := p.Start(); err != nil {
			fmt.Printf("there's been an error: %v", err)
			return
		}
	},
}

func setConfig(config Config) {
	dsn = config.DSN
	afterCompletion = config.AfterCompletion
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pgsp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "display version information")

	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "PostgreSQL data source name")
	rootCmd.PersistentFlags().IntVarP(&afterCompletion, "AfterCompletion", "a", 10, "Number of seconds to display after completion(Seconds)")
	rootCmd.PersistentFlags().Float64VarP(&interval, "Interval", "i", 0.1, "Number of seconds to display after completion(Seconds)")

	_ = viper.BindPFlag("dsn", rootCmd.PersistentFlags().Lookup("dsn"))
	_ = viper.BindPFlag("AfterCompletion", rootCmd.PersistentFlags().Lookup("AfterCompletion"))
	_ = viper.BindPFlag("Interval", rootCmd.PersistentFlags().Lookup("Interval"))
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
		viper.SetConfigName(".pgsp")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal error: %w\n", err)
	}
}
