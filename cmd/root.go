package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"

	"github.com/noborus/pgsp"
	"github.com/noborus/pgsp/tui"

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
	FullScreen      bool    `yaml:"FullScreen"`
}

var (
	verFlag bool
	debug   bool
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
	Version: Version + " rev:" + Revision,
	Run: func(cmd *cobra.Command, args []string) {
		Progress(tui.All)
	},
}

func Progress(targets int) {
	setConfig()
	db, err := pgsp.Connect(config.DSN)
	if err != nil {
		log.Println(err)
		return
	}
	defer pgsp.DisConnect(db)
	if tui.Debug {
		f, err := tea.LogToFile("pgsp.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	model := tui.NewModel(db)
	tui.Targets(&model, targets)

	p := tui.NewProgram(model, config.FullScreen)
	tui.DebugLog("Start")
	if err := p.Start(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		return
	}
	tui.DebugLog("End")
}

func setConfig() {
	tui.AfterCompletion = time.Duration(config.AfterCompletion)
	tui.UpdateInterval = time.Duration(time.Millisecond * time.Duration(config.Interval*1000))
	tui.Debug = debug
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

	rootCmd.PersistentFlags().BoolVarP(&verFlag, "version", "v", false, "display version information")

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug message for toggle")

	var dsn string
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "PostgreSQL data source name")
	_ = viper.BindPFlag("dsn", rootCmd.PersistentFlags().Lookup("dsn"))

	var afterCompletion int
	rootCmd.PersistentFlags().IntVarP(&afterCompletion, "AfterCompletion", "a", 10, "Time to display after completion(Seconds)")
	_ = viper.BindPFlag("AfterCompletion", rootCmd.PersistentFlags().Lookup("AfterCompletion"))

	var interval float64
	rootCmd.PersistentFlags().Float64VarP(&interval, "Interval", "i", 0.5, "Update interval(Seconds)")
	_ = viper.BindPFlag("Interval", rootCmd.PersistentFlags().Lookup("Interval"))

	var fullscreen bool
	rootCmd.PersistentFlags().BoolVarP(&fullscreen, "fullscreen", "f", false, "Display in Full Screen")
	_ = viper.BindPFlag("FullScreen", rootCmd.PersistentFlags().Lookup("fullscreen"))
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

		// Search config in home directory with name ".pgsp" (without extension).
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
