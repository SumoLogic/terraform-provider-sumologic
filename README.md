<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terrafpr," align="right" height="50" />
</a>

# Terraform Provider for Sumo Logic

- Website: [terraform.io](https://terraform.io)
- Chat: [gitter](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing List: [Google Groups](http://groups.google.com/group/terraform-tool)

The Terraform Sumo Logic provider is a plugin for Terraform that allows for the full lifecycle management of Sumo Logic resources.

This provider is maintained by Sumo Logic.

## Getting started

[Using the provider](https://www.terraform.io/docs/providers/sumologic/)

Run `terraform init` to automatically install the latest version of the provider.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x or 0.12.x
- [Go](https://golang.org/doc/install) >= 1.13 (to build the provider plugin)
- [Sumo Logic](https://www.sumologic.com/)

## Developing the provider

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

## Testing the provider

In order to test the provider, you can run `make test`.

For manual testing, run `make install` in the root direcory to install it as a plugin. 
Then run `terraform init` to initialize it.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests *create real resources*, and often cost money to run. The environment variables `SUMOLOGIC_ACCESSID`, `SUMOLOGIC_ACCESSKEY`, and `SUMOLOGIC_ENVIRONMENT` must also be set for acceptance tests to work properly.

[0]: https://help.sumologic.com/Manage/Security/Access-Keys
[1]: https://help.sumologic.com/APIs/General_API_Information/Sumo_Logic_Endpoints_and_Firewall_Security
[10]: https://www.terraform.io/docs/providers/sumologic/
