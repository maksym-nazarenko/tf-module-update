locals {
  region = "eu-west-1"
}

provider "aws" {
  region = local.region
}

module "azure-vpc" {
  source = "git::https://github.com/username/terraform-modules/azure.git//vpc?ref=v0.2.2"
}

module "aws-vpc" {
  source = "git::https://github.com/terraform-aws-modules/terraform-aws.git//src/vpc?ref=v2.78.0"

  name = "simple-example"
  cidr = "10.0.0.0/16"

  azs             = ["${local.region}a", "${local.region}b", "${local.region}c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_ipv6 = false

  public_subnet_tags = {
    Name = "overridden-name-public"
  }

  tags = {
    Owner       = "user"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "vpc-name"
  }
}

module "aws-alerts" {
  source = "git::https://github.com/terraform-aws-modules/alerts.git?ref=v1.0.0"
}
