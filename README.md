<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" alt="Terraform logo" title="Terrafpr," align="right" height="50" />
</a>

# Terraform Provider for Sumo Logic

- Website: [terraform.io](https://terraform.io)
- Chat: [gitter](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing List: [Google Groups](http://groups.google.com/group/terraform-tool)

The Terraform Sumo Logic provider is a plugin for Terraform that allows for the full lifecycle management of Sumo Logic resources.

This provider is maintained by Sumo Logic.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13
- [Go](https://golang.org/doc/install) >= 1.13 (to build the provider plugin)
  - Set [`$GOPATH`](http://golang.org/doc/code.html#GOPATH)
  - Add `$GOPATH/bin` to your `$PATH`
- [Sumo Logic](https://www.sumologic.com/)

## Getting started

Add the Sumo Logic provider to your terraform configuration:

```
terraform {
    required_providers {
        sumologic = {
            source = "sumologic/sumologic"
            version = "" # set the Sumo Logic Terraform Provider version
        }
    }
}
```

Run `terraform init` to automatically install the selected version of the provider.

See the [provider documentation](https://www.terraform.io/docs/providers/sumologic/) for information on the supported resources and example usage.

## Developing the provider

### Build from source

```sh
$ git clone https://github.com/SumoLogic/terraform-provider-sumologic.git
$ cd terraform-provider-sumologic
$ make build
```

This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

### Run locally-built provider

- Run `make install-dev`
- Update the provider source in your configuration file:

```
terraform {
  required_providers {
    sumologic = {
      source = "sumologic.com/dev/sumologic"
      version = "~> 1.0.0"
    }
  }
}
```

## Acceptance tests

> **Note:** Acceptance tests create real resources, and often cost money to run.

### Configuration

Create a personal access key for your Sumo Logic account, and set the following environment variables.

```sh
$ export SUMOLOGIC_ACCESSID="yourAccessID"
$ export SUMOLOGIC_ACCESSKEY="yourAccessKey"
$ export SUMOLOGIC_ENVIRONMENT="yourEnvironment"
$ export SUMOLOGIC_BASE_URL="yourEnvironmentAPIUrl"  # not required for most production deployments
$ export TF_ACC=1
```

More information on configuration can be found [here](https://github.com/SumoLogic/terraform-provider-sumologic/blob/master/website/docs/index.html.markdown#environment-variables).

### Run the tests

```sh
# Run all acceptance tests:
$ make testacc
# Run a specific test:
$ go test -v ./sumologic  -run YourSpecificTestName
```

Some tests require additional configuration for interacting with resources external to Sumo Logic:

- GCP metrics
  - `export SUMOLOGIC_TEST_GOOGLE_APPLICATION_CREDENTIALS=$(cat /path/to/service_acccount.json)`
  - `export SUMOLOGIC_ENABLE_GCP_METRICS_ACC_TESTS="false"` to disable acceptance tests
