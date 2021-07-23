module "azure-db" {
  source = "example.com/terraform-modules/azure.git//azure/db?ref=v1.2.0"

  name = "simple-example"
  cidr = "10.0.0.0/16"
}

module "azure-metric" {
  source = "git::https://example.com/terraform-modules/azure.git//azure/metric?ref=v2.3.0"

  name      = "metric"
  pagerduty = true
}
