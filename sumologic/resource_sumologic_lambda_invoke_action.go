package sumologic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicLambdaInvokeAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSumologicLambdaInvokeActionCreate,
		ReadContext:   resourceSumologicLambdaInvokeActionRead,
		UpdateContext: resourceSumologicLambdaInvokeActionUpdate,
		DeleteContext: resourceSumologicLambdaInvokeActionDelete,

		Schema: map[string]*schema.Schema{
			"lambda_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the AWS Lambda function to invoke.",
			},
			"aws_resource": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN or identifier of the AWS resource to enable logging for.",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The S3 bucket name where logs will be delivered.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A filter expression to match specific resources for enabling logging.",
			},
			"bucket_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The prefix path within the S3 bucket for log delivery.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS account ID where the resources reside.",
			},
			"remove_on_delete_stack": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to remove the logging configuration when the resource is destroyed.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS region where the Lambda function is deployed. If not set, uses AWS_REGION env var or SDK defaults.",
			},
			"last_lambda_output": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The output returned by the most recent Lambda invocation.",
			},
			"last_resource_properties": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON string of the resource properties sent in the most recent Lambda invocation.",
			},
		},
	}
}

var newLambdaClientFunc = defaultNewLambdaClient

func defaultNewLambdaClient(ctx context.Context, region string) (*lambda.Client, error) {
	var opts []func(*awsconfig.LoadOptions) error
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config from environment: %w", err)
	}
	return lambda.NewFromConfig(cfg), nil
}

func resourceSumologicLambdaInvokeActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := newLambdaClientFunc(ctx, d.Get("region").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	payloadMap := map[string]interface{}{
		"action":       "create",
		"ResourceType": "EnableS3LogsResources",
		"ResourceProperties": map[string]interface{}{
			"AWSResource":         d.Get("aws_resource").(string),
			"BucketName":          d.Get("bucket_name").(string),
			"Filter":              d.Get("filter").(string),
			"BucketPrefix":        d.Get("bucket_prefix").(string),
			"AccountID":           d.Get("account_id").(string),
			"RemoveOnDeleteStack": d.Get("remove_on_delete_stack").(bool),
		},
		"OldResourceProperties": map[string]interface{}{},
	}

	output, err := invokeLambda(ctx, client, d.Get("lambda_name").(string), payloadMap)
	if err != nil {
		return diag.FromErr(fmt.Errorf("lambda invocation failed: %w", err))
	}

	d.SetId(fmt.Sprintf("%s-%s-%s", d.Get("lambda_name").(string), d.Get("account_id").(string), d.Get("aws_resource").(string)))
	d.Set("last_lambda_output", output)

	resourceProps, _ := json.Marshal(payloadMap["ResourceProperties"])
	d.Set("last_resource_properties", string(resourceProps))

	return nil
}

func resourceSumologicLambdaInvokeActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSumologicLambdaInvokeActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := newLambdaClientFunc(ctx, d.Get("region").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var oldResourceProperties map[string]interface{}
	if v, ok := d.GetOk("last_resource_properties"); ok {
		if err := json.Unmarshal([]byte(v.(string)), &oldResourceProperties); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse old resource properties: %w", err))
		}
	}

	payloadMap := map[string]interface{}{
		"action":       "update",
		"ResourceType": "EnableS3LogsResources",
		"ResourceProperties": map[string]interface{}{
			"AWSResource":         d.Get("aws_resource").(string),
			"BucketName":          d.Get("bucket_name").(string),
			"Filter":              d.Get("filter").(string),
			"BucketPrefix":        d.Get("bucket_prefix").(string),
			"AccountID":           d.Get("account_id").(string),
			"RemoveOnDeleteStack": d.Get("remove_on_delete_stack").(bool),
		},
		"OldResourceProperties": oldResourceProperties,
	}

	output, err := invokeLambda(ctx, client, d.Get("lambda_name").(string), payloadMap)
	if err != nil {
		return diag.FromErr(fmt.Errorf("lambda invocation failed: %w", err))
	}

	d.Set("last_lambda_output", output)

	resourceProps, _ := json.Marshal(payloadMap["ResourceProperties"])
	d.Set("last_resource_properties", string(resourceProps))

	return nil
}

func resourceSumologicLambdaInvokeActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := newLambdaClientFunc(ctx, d.Get("region").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var resourceProperties map[string]interface{}
	if v, ok := d.GetOk("last_resource_properties"); ok {
		if err := json.Unmarshal([]byte(v.(string)), &resourceProperties); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse resource properties: %w", err))
		}
	}

	payloadMap := map[string]interface{}{
		"action":                "delete",
		"ResourceType":          "EnableS3LogsResources",
		"ResourceProperties":    resourceProperties,
		"OldResourceProperties": map[string]interface{}{},
	}

	_, err = invokeLambda(ctx, client, d.Get("lambda_name").(string), payloadMap)
	if err != nil {
		return diag.FromErr(fmt.Errorf("lambda invocation failed on delete: %w", err))
	}

	d.SetId("")
	return nil
}

func invokeLambda(ctx context.Context, client *lambda.Client, functionName string, payload map[string]interface{}) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      payloadBytes,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke Lambda: %w", err)
	}

	return string(resp.Payload), nil
}
