package cmd

import (
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/ci/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	plan      bool
	deploy    bool
	destroy   bool
	component string

	// arguments that'll be mapped to environment variables.
	argTfStateBucketRegion string
	argTfStateBucket       string
	argTfStateLockTable    string
	argTfRegistryGitHubOrg string
	argTfRegistryBaseURL   string
	argTfVarAWSRegion      string
	argTfVarAWSenvironment string
)

var InfraCMD = &cobra.Command{
	Use: "infra",
	Long: `Build and deploy the infrastructure.
This command will build and deploy the underlying infrastructure for this Rotator lambda.`,
	Example: `
rotator infra deploy
  `,
	Run: func(cmd *cobra.Command, args []string) {
		if plan {
			err := tasks.InfraPlan()
			if err != nil {
				os.Exit(1)
			}
		}

		if deploy {
			err := tasks.InfraDeploy()
			if err != nil {
				os.Exit(1)
			}
		}

		if destroy {
			err := tasks.InfraDestroy()
			if err != nil {
				os.Exit(1)
			}
		}
		//_ = cmd.Help()
		os.Exit(0)
	},
}

func addInfraCMDFlags() {
	InfraCMD.Flags().BoolVarP(&plan, "plan", "p", false, "Plan the infrastructure deployment.")
	InfraCMD.Flags().BoolVarP(&deploy, "deploy", "k", false, "Deploy the infrastructure.")
	InfraCMD.Flags().BoolVarP(&destroy, "destroy", "", false, "Destroy the infrastructure.")
	InfraCMD.Flags().StringVarP(&component, "component", "c", "", "The component to deploy.")

	// arguments that'll be mapped to environment variables.
	InfraCMD.Flags().StringVarP(&argTfStateBucketRegion, "tf-state-bucket-region", "", "", "The region of the S3 bucket to store the Terraform state.")
	InfraCMD.Flags().StringVarP(&argTfStateBucket, "tf-state-bucket", "", "", "The name of the S3 bucket to store the Terraform state.")
	InfraCMD.Flags().StringVarP(&argTfStateLockTable, "tf-state-lock-table", "", "", "The name of the DynamoDB table to lock the Terraform state.")
	InfraCMD.Flags().StringVarP(&argTfRegistryGitHubOrg, "tf-registry-github-org", "", "", "The GitHub organisation to use for the Terraform registry.")
	InfraCMD.Flags().StringVarP(&argTfRegistryBaseURL, "tf-registry-base-url", "", "", "The base URL to use for the Terraform registry.")
	InfraCMD.Flags().StringVarP(&argTfVarAWSRegion, "tf-var-aws-region", "", "", "The AWS region to deploy to.")
	InfraCMD.Flags().StringVarP(&argTfVarAWSenvironment, "tf-var-aws-environment", "", "", "The AWS environment to deploy to.")

	_ = viper.BindPFlag("plan", InfraCMD.Flags().Lookup("plan"))
	_ = viper.BindPFlag("deploy", InfraCMD.Flags().Lookup("deploy"))
	_ = viper.BindPFlag("destroy", InfraCMD.Flags().Lookup("destroy"))
	_ = viper.BindPFlag("component", InfraCMD.Flags().Lookup("component"))

	// arguments that'll be mapped to environment variables.
	_ = viper.BindPFlag("tf-state-bucket-region", InfraCMD.Flags().Lookup("tf-state-bucket-region"))
	_ = viper.BindPFlag("tf-state-bucket", InfraCMD.Flags().Lookup("tf-state-bucket"))
	_ = viper.BindPFlag("tf-state-lock-table", InfraCMD.Flags().Lookup("tf-state-lock-table"))
	_ = viper.BindPFlag("tf-registry-github-org", InfraCMD.Flags().Lookup("tf-registry-github-org"))
	_ = viper.BindPFlag("tf-registry-base-url", InfraCMD.Flags().Lookup("tf-registry-base-url"))
	_ = viper.BindPFlag("tf-var-aws-region", InfraCMD.Flags().Lookup("tf-var-aws-region"))
	_ = viper.BindPFlag("tf-var-aws-environment", InfraCMD.Flags().Lookup("tf-var-aws-environment"))
}

func init() {
	addInfraCMDFlags()
}
