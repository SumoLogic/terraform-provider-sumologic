[![Build Status](https://travis-ci.org/sumologic/sumologic-terraform-provider.svg?branch=master)](https://travis-ci.org/sumologic/sumologic-terraform-provider) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/sumologic/sumologic-terraform-provider/issues)

# terraform-provider-sumologic
This provider is used to manage multiple configuration entities within the Sumo Logic product.

## Support

The code in this repository has been developed in collaboration with the Sumo Logic community and is not supported via standard Sumo Logic Support channels. For any issues or questions please submit an issue within the GitHub repository. The maintainers of this project will work directly with the community to answer any questions, address bugs, or review any requests for new features. 

## License
Released under Apache 2.0 License.

# Getting started / usage

See [docs][10]

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
[10]: r/docs
