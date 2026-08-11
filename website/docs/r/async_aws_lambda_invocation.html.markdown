---
layout: "sumologic"
page_title: "SumoLogic: sumologic_async_aws_lambda_invocation"
description: |-
  Provides a Sumologic Async AWS Lambda Invocation resource
---

# sumologic_async_aws_lambda_invocation

Provides a resource to asynchronously invoke an AWS Lambda function. Uses Lambda's `Event` invocation type, so the function is triggered and Terraform does not wait for it to complete. Useful for triggering background setup operations during infrastructure provisioning.

This resource creates its own AWS Lambda client. It does not use the Sumo Logic provider credentials for AWS operations.

## Example Usage

```hcl
resource "sumologic_async_aws_lambda_invocation" "trigger_setup" {
  function_name = "MySetupFunction"
  region        = "us-east-1"
  aws_profile   = "my-aws-profile"
  input         = jsonencode({ env = "production" })
  qualifier     = "$LATEST"
}
```

### Multi-profile example

```hcl
resource "sumologic_async_aws_lambda_invocation" "cross_account" {
  function_name = "SetupFunction"
  region        = "us-west-2"
  aws_profile   = "prod-account"
  input         = jsonencode({ source = "terraform" })
}
```

## Argument Reference

The following arguments are supported:

- `function_name` - (Required, Forces new resource) The name or ARN of the AWS Lambda function to invoke.
- `region` - (Required, Forces new resource) The AWS region where the Lambda function is deployed.
- `input` - (Optional, Forces new resource) JSON string payload to pass to the Lambda function. Defaults to `"{}"`.
- `qualifier` - (Optional, Forces new resource) The Lambda function version or alias to invoke. Defaults to `"$LATEST"`.
- `aws_profile` - (Optional, Forces new resource) The AWS credentials profile to use. If not set, uses the default credential chain (`AWS_PROFILE` env var, instance profile, etc.).
- `triggers` - (Optional, Forces new resource) A map of arbitrary key/value pairs. Changing any value forces re-invocation, useful for triggering a new invocation without changing other arguments.

## Attributes Reference

The following attributes are exported:

- `id` - Identifier in the format `<function_name>-<status_code>`.
- `status_code` - HTTP status code returned by the Lambda invocation API (`202` indicates the async invocation was accepted).

## AWS Authentication

AWS credentials are resolved in the following order:

1. `aws_profile` argument (if set)
2. `AWS_PROFILE` environment variable
3. `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` environment variables
4. Instance profile / ECS task role / Web Identity token

## Import

This resource does not support import.
