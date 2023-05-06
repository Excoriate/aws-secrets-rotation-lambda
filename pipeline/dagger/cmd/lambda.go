package cmd

import (
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/ci/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	packageZip bool
	compile    bool
	lambdaSrc  string
)

var LambdaCMD = &cobra.Command{
	Use: "lambda",
	Long: `Perform actions such as compress the lambda source code to package it into a
deployable artifact, among others`,
	Example: `
rotator lambda package-zip
  `,
	Run: func(cmd *cobra.Command, args []string) {
		if compile {
			err := tasks.CompileLambda()
			if err != nil {
				os.Exit(1)
			}
		}

		if packageZip {
			err := tasks.PackageZip()
			if err != nil {
				os.Exit(1)
			}
		}
		//_ = cmd.Help()
		os.Exit(0)
	},
}

func addFlags() {
	LambdaCMD.Flags().BoolVarP(&packageZip, "package-zip", "z", false,
		"Generate a zip file with the lambda code.")
	LambdaCMD.Flags().BoolVarP(&compile, "compile", "c", false, "Compile the lambda source code.")
	LambdaCMD.Flags().StringVarP(&lambdaSrc, "lambda-src", "s", ".",
		"Lambda source code directory. If it's not set, it'll use the current directory.")
	_ = viper.BindPFlag("package-zip", LambdaCMD.Flags().Lookup("generate-zip"))
	_ = viper.BindPFlag("lambda-src", LambdaCMD.Flags().Lookup("lambda-src"))
	_ = viper.BindPFlag("compile", LambdaCMD.Flags().Lookup("compile"))
}

func init() {
	addFlags()
}
