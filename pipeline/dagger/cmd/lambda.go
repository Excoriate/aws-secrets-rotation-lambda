package cmd

import (
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/ci/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	generateZip bool
	lambdaSrc   string
)

var LambdaCMD = &cobra.Command{
	Use: "lambda",
	Long: `Perform actions such as compress the lambda source code to package it into a
deployable artifact, among others`,
	Example: `
rotator lambda generate-zip
  `,
	Run: func(cmd *cobra.Command, args []string) {
		if generateZip {
			err := tasks.GenerateZip()
			if err != nil {
				os.Exit(1)
			}
		}
		//_ = cmd.Help()
		os.Exit(0)
	},
}

func addFlags() {
	LambdaCMD.Flags().BoolVarP(&generateZip, "generate-zip", "z", false, "Generate a zip file with the lambda code.")
	LambdaCMD.Flags().StringVarP(&lambdaSrc, "lambda-src", "s", ".",
		"Lambda source code directory. If it's not set, it'll use the current directory.")
	_ = viper.BindPFlag("generate-zip", LambdaCMD.Flags().Lookup("generate-zip"))
	_ = viper.BindPFlag("lambda-src", LambdaCMD.Flags().Lookup("lambda-src"))
}

func init() {
	addFlags()
}
