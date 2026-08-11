package sumologic

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var newAsyncLambdaClientFunc = defaultNewAsyncLambdaClient

func defaultNewAsyncLambdaClient(ctx context.Context, region string, profile string) (*lambda.Client, error) {
	var opts []func(*awsconfig.LoadOptions) error
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}
	if profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}
	return lambda.NewFromConfig(cfg), nil
}

func resourceSumologicAsyncAwsLambdaInvocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsyncLambdaInvocationCreate,
		ReadContext:   resourceAsyncLambdaInvocationRead,
		DeleteContext: resourceAsyncLambdaInvocationDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"input": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "{}",
				ValidateFunc: validation.StringIsJSON,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "$LATEST",
			},
			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status_code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aws_profile": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
		},
	}
}

func resourceAsyncLambdaInvocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	region := d.Get("region").(string)
	functionName := d.Get("function_name").(string)
	input := d.Get("input").(string)
	qualifier := d.Get("qualifier").(string)

	client, err := newAsyncLambdaClientFunc(ctx, region, d.Get("aws_profile").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        []byte(input),
		Qualifier:      aws.String(qualifier),
	})
	if err != nil {
		return diag.Errorf("async Lambda invocation failed for %s: %s", functionName, err)
	}

	d.SetId(fmt.Sprintf("%s-%d", functionName, resp.StatusCode))
	d.Set("status_code", int(resp.StatusCode))
	return nil
}

func resourceAsyncLambdaInvocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceAsyncLambdaInvocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
