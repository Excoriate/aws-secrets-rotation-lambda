include "root" {
  path           = find_in_parent_folders()
  merge_strategy = "deep"
  expose        = true
}

include "registry" {
  path           = "${get_terragrunt_dir()}/../_env/_registry.hcl"
  expose         = true
  merge_strategy = "deep"
}

include "metadata" {
  path           = "${get_terragrunt_dir()}/../_env/_metadata.hcl"
  expose         = true
  merge_strategy = "deep"
}

locals {
  /*
    NOTE:
      ----------------------------------------------------------
      * Customise registry values accordingly.
      * These are the values that will be used to construct the module source URL.
      ----------------------------------------------------------
  */
  module_repo =  get_env("TF_MODULE_REPO", "terraform-registry-aws-events")
  module_path = get_env("TF_MODULE_PATH", "modules/lambda/lambda-function")
  module_version = get_env("TF_MODULE_VERSION", "v0.1.12")
  registry_base_url = include.registry.locals.registry_base_url
  registry_github_org = include.registry.locals.registry_github_org
  /*
   NOTE:
     ----------------------------------------------------------
     * Customise tags accordingly.
     * These 'tags' values are merged with what's defined in the
      '_metadata.hcl' file.
     ----------------------------------------------------------
  */
  tags = {}
  source_url = format("%s/%s/%s//%s?ref=%s", local.registry_base_url, local.registry_github_org, local.module_repo, local.module_path, local.module_version)
  source_url_show = run_cmd("sh", "-c", format("export SOURCE_URL=%s; echo source url : [$SOURCE_URL]", local.source_url))

  lambda_name = format("%s-%s-secrets-manager-rotator-function-%s", get_env("TF_VAR_environment", "dev"), get_env("TF_VAR_aws_region"), get_env("TF_VAR_rotator_lambda_name"))
  deployment_bucket = format("%s-%s-secrets-manager-rotator-deployments-%s", get_env("TF_VAR_environment", "dev"), get_env("TF_VAR_aws_region"), get_env("TF_VAR_rotator_lambda_name"))
}

terraform {
  source = format("%s/%s/%s//%s?ref=%s", local.registry_base_url, local.registry_github_org, local.module_repo, local.module_path, local.module_version)
}

inputs = {
  tags = merge(include.metadata.locals.tags, {
    "Name" = "lambda-rotator-function"
  })

  lambda_config = [
    {
      name = local.lambda_name
      name = local.lambda_name
      handler       = "secrets-manager-rotator-lambda"
      runtime = "go1.x"
      publish = true
      deployment_type = {
        from_s3_existing_file = true
      }
    }
  ]

  lambda_observability_config = [
    {
      name         = local.lambda_name
      logs_enabled = true
      tracing_enabled = true
      tracing_mode = "Active"
    }
  ]

  lambda_permissions_config = [
    {
      name         = local.lambda_name
    }
  ]

  lambda_enable_secrets_manager = [
    {
      name        = local.lambda_name
      secret_name = get_env("TF_VAR_secret_name")
    }
  ]

  lambda_s3_from_existing_config = [
    {
      name              = local.lambda_name
      s3_bucket         = local.deployment_bucket
      s3_key            = "releases/secrets-manager-rotator-lambda.zip"
    }
  ]

  lambda_enable_secrets_manager_rotation = [
    {
      name              = local.lambda_name
      secrets_to_rotate = [get_env("TF_VAR_secret_name")]
    }
  ]

  lambda_host_config = [
    {
      name = local.lambda_name
      environment_variables = {
        "IS_ENABLED"  = get_env("TF_VAR_rotation_lambda_enabled", "true")
      }
    }
  ]

}
