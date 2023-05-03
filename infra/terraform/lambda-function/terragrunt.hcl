include "root" {
  path           = find_in_parent_folders()
  merge_strategy = "deep"
  expose        = true
}

terraform {
  source = "git::https://github.com/wbd-streaming/infra-global-manifest-stitcher-service//terraform/modules/irsa?depth=1&ref=v0.1.3"
}

inputs = {
  tags = {}
}
