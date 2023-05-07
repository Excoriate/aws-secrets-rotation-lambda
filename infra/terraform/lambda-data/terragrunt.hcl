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

terraform {
  source = "./"
}

inputs = {
  secret_name = get_env("TF_VAR_secret_name")
}
