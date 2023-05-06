package cmd

import (
	"context"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	dryRun  bool
	debug   bool
	cfgFile string
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "rotator",
	Long: `Rotator is a command-line tool that helps you to rotate your AWS credentials,
powered by Dagger.io.`,
	Example: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Show main pipeline title
		uxTitle := tui.NewTitle()
		uxTitle.ShowTitle("rotator-lambda")

		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func addPersistentFlags() {
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode. If so, it'll not perform any changes.")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Debug mode. If so, it'll print debug messages.")
	_ = viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rotator-cli")

		_ = viper.SafeWriteConfig()
		//if err != nil {
		//	// Check if error relates to the file already exist.
		//	// If it does, then it's fine, otherwise, exit.
		//	if !os.IsExist(err) {
		//		fmt.Println(err)
		//		os.Exit(1)
		//	}
		//	//os.Exit(1)
		//}
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func AddSubCommands() {
	rootCmd.AddCommand(LambdaCMD)
	rootCmd.AddCommand(InfraCMD)
}

func init() {
	cobra.OnInitialize(initConfig)
	addPersistentFlags()

	AddSubCommands()
}
