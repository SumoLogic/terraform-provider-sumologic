[![Build Status](https://travis-ci.org/SumoLogic/sumologic-terraform-provider.svg?branch=master)](https://travis-ci.org/SumoLogic/sumologic-terraform-provider) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/SumoLogic/sumologic-terraform-provider/issues)

# terraform-provider-sumologic
This provider is used to manage Hosted collectors and sources supported by Sumo Logic.

## Support

The code in this repository has been developed in collaboration with the Sumo Logic community and is not supported via standard Sumo Logic Support channels. For any issues or questions please submit an issue within the GitHub repository. The maintainers of this project will work directly with the community to answer any questions, address bugs, or review any requests for new features. 

## License
Released under Apache 2.0 License.

## Usage

The provider needs to be configured with the proper credentials before it can be used.  You must provide an [Access ID and Access Key][0] to use this provider.

### Authentication
The Sumo Logic Provider offers a flexible means of providing credentials for authentication. The following methods are supported and explained below:

 - Static credentials
 - Environment variables

#### Static credentials
Static credentials can be provided by adding an `access_id` and `access_key` in-line in the Sumo Logic provider block:

Usage:
```hcl
provider "sumologic" {
    access_id   = "your-access-id"
    access_key  = "your-access-key"
}
```

#### Environment variables
You can provide your credentials via the `SUMOLOGIC_ACCESSID` and `SUMOLOGIC_ACCESSKEY` environment variables, representing your Sumo Logic Access ID and Sumo Logic Access Key, respectively.

Usage:
```hcl
provider "sumologic" { }
```

```bash
$ export SUMOLOGIC_ACCESSID="your-access-id"
$ export SUMOLOGIC_ACCESSKEY="your-access-id"
$ terraform plan
```

## Argument Reference
- `access_id` - (Optional) This is the Sumo Logic Access ID. It must be provided, but it can also be sourced from the SUMOLOGIC_ACCESSID environment variable.
- `access_key` - (Optional) This is the Sumo Logic Access Key. It must be provided, but it can also be sourced from the SUMOLOGIC_ACCESSKEY environment variable.
- `environment` - (Optional) This is the API endpoint to use. Default is us2, but it can also be sourced from the SUMOLOGIC_ENVIRONMENT environment variable. See the [Sumo Logic documentation][1] for details on which environment you should use.

# Building the provider

In this section you will learn how to build and run terraform-provider-sumologic locally. Please follow the steps below:

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) >= 1.9 (to build the provider plugin)
- [Sumo Logic](https://www.sumologic.com/pricing/)

Build and install
------------
`make install`

[0]: https://help.sumologic.com/Manage/Security/Access-Keys
[1]: https://help.sumologic.com/APIs/General_API_Information/Sumo_Logic_Endpoints_and_Firewall_Security
