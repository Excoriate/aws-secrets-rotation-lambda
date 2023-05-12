<h1 align="center">
  AWS Secrets Rotator ⚙️️
</h1>
<p align="center">Simple, and <b>fully-automated</b> AWS Lambda (Secret 🔑) rotator! <b> function, that works on top of very nice things! ❤️️</b>.<br/><br/>

---
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md) [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)[![Setup Demo Config](https://github.com/Excoriate/aws-secrets-rotation-lambda/actions/workflows/setup-demo.yaml/badge.svg)](https://github.com/Excoriate/aws-secrets-rotation-lambda/actions/workflows/setup-demo.yaml)[![release](https://github.com/Excoriate/aws-secrets-rotation-lambda/actions/workflows/release.yaml/badge.svg)](https://github.com/Excoriate/aws-secrets-rotation-lambda/actions/workflows/release.yaml)


---

## Description

This is a simple project that implements a **rotator** [AWS Lambda](https://aws.amazon.com/lambda/) function, that works on top of [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) to rotate secrets. For more details about how the process works, please refer to the [AWS Secrets Manager documentation](https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets.html). In addition, it includes the entire architecture as infrastructure as code to support this functionality.

## Stack

- Beside the AWS Lambda function, this project also implement the stack required to deploy the function, and the required resources to make it work.
- The stack is implemented using [Terragrunt](https://terragrunt.gruntwork.io/), and [Terraform](https://www.terraform.io/).
- The pipeline is fully portable, and it's using [Dagger.IO](https://dagger.io/) to build and deploy the stack.

## Architecture

### Layers

```mermaid
graph TB
  A[Deployment Bucket (S3)] --> B[Rotation Configuration (Terraform Module)]
  B --> C[Target Secret]
  C --> D[Lambda Function]

```

### Infrastructure/Platform built-in capabilities

These configurations (infrastructure as code) are implemented:

- [x] [AWS Lambda](https://aws.amazon.com/lambda/) rotator function (Golang `1.20`).
- [x] [AWS Lambda](https://aws.amazon.com/lambda/) underlaying infrastructure (role, policies, lambda permissions, etc.). It uses this [Terraform module](https://github.com/Excoriate/terraform-registry-aws-events/tree/main/modules/lambda/lambda-function) to implement the lambda function.
- [x] [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) secret (sample secret). It uses this [Terraform module](https://github.com/Excoriate/terraform-registry-aws-storage/tree/main/modules/secrets-manager) to implement the secret.

>**Note**: Ensure you are going to replace the _sample secret_ in the `infra/terraform/secrets-manager-secret` module, with your own secret.

- [x] [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) rotation configuration. It uses this [Terraform module](https://github.com/Excoriate/terraform-registry-aws-storage/tree/main/modules/secrets-manager-rotation) to implement the rotation configuration.
- [x] [S3 Deployment Bucket](https://aws.amazon.com/s3/) to store the lambda function code. It uses this [Terraform module](https://github.com/Excoriate/terraform-registry-aws-storage/tree/main/modules/s3/s3-lambda-deployment-bucket) to implement the bucket.


## Configuration

Ensure to export this environment variables (see the `.env.example` file):

```bash
TF_STATE_BUCKET_REGION="eu-central-1"
TF_STATE_BUCKET="my_state_bucket"
TF_STATE_LOCK_TABLE="my_state_lock_table"
TF_REGISTRY_GITHUB_ORG="Excoriate"
TF_REGISTRY_BASE_URL="git::https://github.com"
TF_VAR_aws_region="us-east-1"
TF_VAR_environment="dev"
TF_VAR_rotator_lambda_name="demo"
TF_VAR_secret_name="/dev/us-east-1/secrets-manager-rotator-demo/my-demo-secret-to-rotate-1"
TF_VAR_rotation_schedule="cron(0 7/4 * * ? *)"
TF_VAR_rotation_lambda_enabled=true
TF_VERSION="v1.4.6"
TG_VERSION="v0.42.8"
```

This configuration is also available in the two sample **GitHub Actions Workflow** files: [demo-create](.github/workflows/demo-create.yml), and [demo-destroy](.github/workflows/demo-destroy.yml).

### Local execution

This project uses [Dagger.io](https://dagger.io/) pipelines as code, everything that runs in GitHub Actions can run in your local machine. The pipeline provide a set of built-in commands that you can wrapped in a [Taskfile](https://taskfile.dev/#/).
Run the following command to see the available commands:

```bash
task pipeline-dagger-run --
```

#### Examples

```bash
# Make a terragrunt plan on the lambda function module
task pipeline-dagger-run -- infra --component=lambda --plan

# Generate the compiled binary and zip it for the lambda function to be deployed
task pipeline-dagger-run -- lambda --package-zip --lambda-src=src/lambda/secrets-manager-rotator-go

# Upload the packaged lambda to S3
task pipeline-dagger-run -- lambda --upload-to-s3 --s3-bucket=dev-us-east-1-secrets-manager-rotator-deployments-demo \
--s3-destination-path=releases \
--s3-file-to-upload=output/lambda-zip/linux/amd64/secrets-manager-rotator-lambda.zip
```

>**Note**: Ensure that the necessary `AWS_*` environment variables are exported.

## Roadmap 🗓️

There are more things to do, however, the following are the main ones:

- [ ] Add support for DB related secrets.
- [ ] Add/Cover more test cases.

>**Note**: This is still work in progress, however, I'll be happy to receive any feedback or contribution. Ensure you've read the [contributing guide](./CONTRIBUTING.md) before doing so.


## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).

## Community

Find me in:

- 📧 [Email](mailto:alex@ideaup.cl)
- 🧳 [Linkedin](https://www.linkedin.com/in/alextorresruiz/)


<a href="https://github.com/Excoriate/stiletto/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Excoriate/stiletto" />
</a>
