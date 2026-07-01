---
layout: "sumologic"
page_title: "SumoLogic: sumologic_lambda_invoke_action"
description: |-
  Provides a Sumologic Lambda Invoke Action resource
---

# sumologic_lambda_invoke_action
Provides a resource to invoke an AWS Lambda function for enabling S3 logging and other auto-enable operations as part of AWS Observability setup.

This resource creates its own AWS Lambda client using environment variables (`AWS_REGION`, `AWS_PROFILE` or `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`). It does not use the Sumo Logic provider credentials for AWS operations.

## Example Usage

```hcl
resource "sumologic_lambda_invoke_action" "enable_s3_logging" {
  lambda_name            = "SumologicEnableExistingResources"
  aws_resource           = "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/abc123"
  bucket_name            = "my-access-logs-bucket"
  filter                 = "'Type': 'application'|'type': 'application'"
  bucket_prefix          = "elasticloadbalancing"
  account_id             = "123456789012"
  remove_on_delete_stack = true
}
```

## Argument Reference

The following arguments are supported:

- `lambda_name` - (Required) The name of the AWS Lambda function to invoke.
- `aws_resource` - (Required) The ARN or identifier of the AWS resource to enable logging for.
- `bucket_name` - (Required) The S3 bucket name where logs will be delivered.
- `account_id` - (Required) The AWS account ID.
- `filter` - (Optional) A filter expression to match specific resources. Defaults to `""`.
- `bucket_prefix` - (Optional) The prefix path within the S3 bucket for log delivery. Defaults to `""`.
- `remove_on_delete_stack` - (Optional) Whether to remove the logging configuration when the resource is destroyed. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

- `id` - The unique identifier for this resource (format: `<lambda_name>-<timestamp>`).
- `last_lambda_output` - The output returned by the most recent Lambda invocation.
- `last_resource_properties` - JSON string of the resource properties sent in the most recent Lambda invocation.

## AWS Authentication

This resource requires AWS credentials configured via environment variables:

```bash
export AWS_REGION="us-east-1"
export AWS_PROFILE="my-profile"
# or
export AWS_ACCESS_KEY_ID="..."
export AWS_SECRET_ACCESS_KEY="..."
```

The resource does **not** use the Sumo Logic provider's `access_id`/`access_key` for AWS operations.

## Import

This resource does not support import.
