[![Build Status](https://travis-ci.org/sumologic/sumologic-terraform-provider.svg?branch=master)](https://travis-ci.org/sumologic/sumologic-terraform-provider) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/sumologic/sumologic-terraform-provider/issues)

# terraform-provider-sumologic
This provider is used to manage multiple configuration entities within the Sumo Logic product.

## Support

The code in this repository has been developed in collaboration with the Sumo Logic community and is not supported via standard Sumo Logic Support channels. For any issues or questions please submit an issue within the GitHub repository. The maintainers of this project will work directly with the community to answer any questions, address bugs, or review any requests for new features. 

## License
Released under Mozilla Public License 2.0.

# Getting started / usage

See [docs][10]

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x or 0.12.x
- [Go](https://golang.org/doc/install) >= 1.9 (to build the provider plugin)
- [Sumo Logic](https://www.sumologic.com/pricing/)

# Using the provider

To use the provider run `make install` in the root direcory to install it as a plugin. You can then run `terraform init` to initialize it.

# Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine. You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone repository to: `$GOPATH/src/SumoLogic/sumologic-terraform-provider`

```sh
$ mkdir -p $GOPATH/src/SumoLogic;
$ cd $GOPATH/src/SumoLogic
$ git clone git@github.com:SumoLogic/sumologic-terraform-provider.git
```

Enter the provider directory and build the provider. To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ cd $GOPATH/src/SumoLogic/sumologic-terraform-provider
$ make build
```

# Testing the provider

In order to test the provider, you can run `make test`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests *create real resources*, and often cost money to run. The environment variables `SUMOLOGIC_ACCESSID`, `SUMOLOGIC_ACCESSKEY`, and `SUMOLOGIC_ENVIRONMENT` must also be set for acceptance tests to work properly.

[0]: https://help.sumologic.com/Manage/Security/Access-Keys
[1]: https://help.sumologic.com/APIs/General_API_Information/Sumo_Logic_Endpoints_and_Firewall_Security
[10]: website/docs/README.md
