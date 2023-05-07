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
  module_repo =  get_env("TF_MODULE_REPO", "terraform-registry-aws-storage")
  module_path = get_env("TF_MODULE_PATH", "modules/secrets-manager")
  module_version = get_env("TF_MODULE_VERSION", "v1.2.1")
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
}

terraform {
  source = format("%s/%s/%s//%s?ref=%s", local.registry_base_url, local.registry_github_org, local.module_repo, local.module_path, local.module_version)
}

inputs = {
  tags = merge(include.metadata.locals.tags, {
    "Name" = "lambda-rotator-secrets-manager-demo"
  })
  secrets_config = [
    {
      name = format("%s-secrets-manager-rotator-%s-secret", get_env("TF_VAR_environment", "dev"), get_env("TF_VAR_rotator_lambda_name"))
      path = format("/%s/secrets-manager-rotator/%s/my-demo-secret", get_env("TF_VAR_environment", "dev"), get_env("TF_VAR_rotator_lambda_name"))
      enable_random_secret_value = true
    }
  ]
}