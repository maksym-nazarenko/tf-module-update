Terraform module source updater
===============================

A simple tool to update host, module path or revision of a module.

**Disclaimer**
>This tool is under heavy development, so no part of it may be considered as stable.  
Any commit may break the things so be sure to pin version to particular commit (releases are coming).

> **NOTE**: Only git-over-https is well tested supported at the moment. More schemes and use cases might be added in the future

Assume, there are a lot of blocks like these in your Terraform setup:

*main.tf*
```terraform
module "vpc" {
  source = "git::https://github.com/example-corp/terraform-modules.git//src/vpc?ref=v1.2.3"
}

module "rds" {
  source = "git::https://github.com/example-corp/terraform-modules.git//src/rds?ref=v1.2.9"
}
```

*app/backend/main.tf*
```terraform
module "queue" {
  source = "git::https://github.com/example-corp/terraform-modules.git//src/sqs?ref=v1.1.3"
}

module "lambda" {
  source = "git::https://github.com/example-corp/terraform-modules.git//src/lambda?ref=v1.1.3"
}
```

*mobile/backend/main.tf*
```terraform
module "api-gw" {
  source = "git::https://github.com/example-corp/terraform-modules.git//src/apigw?ref=v1.8.3"
}
```

so you end up with bunch of versions of the same repository.

This tool allows you to update desired module to a defined version in seconds.

Or have it as part of CI to be more declarative.

## Usage

### As CLI binary
```shell
$ tf-module-update -from.url='git::https://github.com/example-corp/terraform-modules.git//src/apigw?ref=v1.8.3' -to.revision='v2.1.1' /path/to/terraform/files /another/path/goes/here
```

CLI flags

There are 2 main flag sets for module source filtering and manipulation:
- `'-from.*'` performs filtering of modules in `*.tf` files that will be considered for updating
- `'-to.*'` builds patches to apply to source URL of filtered module

|Flag|Meaning|Example|
|----|-------|-------|
|`*.url`|The full url of module source|https://github.com/example-org/tf-modules.git//aws/vpc/multizone?ref=v1.0.0|
|`*.scheme`|Source scheme|`https`, `http`|
|`*.host`|Host part|github.com|
|`*.module`|Module part|/example-org/tf-modules/aws/vpc|
|`*.submodule`|Submodule, if source scheme supports it|//src/subfolder/azure|
|`*.revision`|Revision, usually a tag in form of `vX.Y.Z`|v2.0.5|
|`-write`|Boolean flag to perform actual update on original files. `Default` (not set) is `false`||
|`-log.level`|Level of logging for application. `Default` is `info`|-log.level=debug|

Both, `from.*` and `to.*` flag sets have only one rule for ordering: `*.url`, if present, builds the initial object and specific flags like `*.submodule` update it.

In this example:
```shell
$ tf-module-update -from.url='https://github.com/example-org/tf-modules.git//aws/vpc/multizone?ref=v1.0.0' -from.submodule='//aws/vpc/subnets' -to.revision='v1.2.1'
```
initial object to filter module sources will be built using full `-from.url` string but submodule is changed to `//aws/vpc/subnets`.

It is equivalent to this call:
```shell
$ tf-module-update -from.scheme='https' -from.host='github.com' -from.module='/example-org/tf-modules.git' -from.submodule='//aws/vpc/subnets' -from.revision='v1.0.0' -to.revision='v1.2.1'
```
but the previous example is less verbose.

The resulting logic would be: replace all occurrences of `https://github.com/example-org/tf-modules.git//aws/vpc/subnets?ref=v1.0.0` with `https://github.com/example-org/tf-modules.git//aws/vpc/subnets?ref=v1.2.1`


### As package in another project

`TBD`: pull the code out of `internal` folder


### Examples

You can find a few `.tf` files to play with under [examples/fixture](examples/fixtures) directory

`TBD`

### TBDs

* add `Development` section
* add Docker environment
* add Github Actions
  * deploy Docker images
* make code available as package
