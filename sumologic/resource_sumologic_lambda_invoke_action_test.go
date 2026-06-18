package sumologic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newMockLambdaClient(handler http.HandlerFunc) *lambda.Client {
	server := httptest.NewServer(handler)
	client := lambda.New(lambda.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(server.URL),
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Source:          "test",
			}, nil
		}),
	})
	return client
}

func withMockLambdaClient(handler http.HandlerFunc, fn func()) {
	client := newMockLambdaClient(handler)
	original := newLambdaClientFunc
	newLambdaClientFunc = func(ctx context.Context, region string) (*lambda.Client, error) {
		return client, nil
	}
	defer func() { newLambdaClientFunc = original }()
	fn()
}

func lambdaSuccessHandler(response string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"StatusCode": 200, "Payload": %q}`, response)))
	}
}

func lambdaErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
	}
}

func newLambdaTestResourceData() *schema.ResourceData {
	r := resourceSumologicLambdaInvokeAction()
	d := r.TestResourceData()
	d.Set("lambda_name", "my-lambda")
	d.Set("aws_resource", "my-resource")
	d.Set("bucket_name", "my-bucket")
	d.Set("filter", "*.log")
	d.Set("bucket_prefix", "logs/")
	d.Set("account_id", "123456789012")
	d.Set("remove_on_delete_stack", true)
	return d
}

// --- Schema Tests ---

func TestLambdaInvokeResourceSchema(t *testing.T) {
	r := resourceSumologicLambdaInvokeAction()

	expectedAttrs := map[string]struct {
		typ      schema.ValueType
		required bool
		optional bool
		computed bool
	}{
		"lambda_name":              {typ: schema.TypeString, required: true},
		"aws_resource":             {typ: schema.TypeString, required: true},
		"bucket_name":              {typ: schema.TypeString, required: true},
		"filter":                   {typ: schema.TypeString, optional: true},
		"bucket_prefix":            {typ: schema.TypeString, optional: true},
		"account_id":               {typ: schema.TypeString, required: true},
		"remove_on_delete_stack":   {typ: schema.TypeBool, optional: true},
		"region":                   {typ: schema.TypeString, optional: true},
		"last_lambda_output":       {typ: schema.TypeString, computed: true},
		"last_resource_properties": {typ: schema.TypeString, computed: true},
	}

	for name, expected := range expectedAttrs {
		attr, ok := r.Schema[name]
		if !ok {
			t.Errorf("schema missing attribute '%s'", name)
			continue
		}
		if attr.Type != expected.typ {
			t.Errorf("attribute '%s': expected type %v, got %v", name, expected.typ, attr.Type)
		}
		if attr.Required != expected.required {
			t.Errorf("attribute '%s': expected Required=%v, got %v", name, expected.required, attr.Required)
		}
		if attr.Optional != expected.optional {
			t.Errorf("attribute '%s': expected Optional=%v, got %v", name, expected.optional, attr.Optional)
		}
		if attr.Computed != expected.computed {
			t.Errorf("attribute '%s': expected Computed=%v, got %v", name, expected.computed, attr.Computed)
		}
	}
}

// --- InvokeLambda Tests ---

func TestInvokeLambda_Success(t *testing.T) {
	client := newMockLambdaClient(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"StatusCode": 200,
			"Payload":    `{"status": "success"}`,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	payload := map[string]interface{}{
		"action":       "create",
		"ResourceType": "EnableS3LogsResources",
		"ResourceProperties": map[string]interface{}{
			"BucketName": "test-bucket",
		},
	}

	output, err := invokeLambda(context.Background(), client, "test-function", payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestInvokeLambda_Error(t *testing.T) {
	client := newMockLambdaClient(lambdaErrorHandler())

	payload := map[string]interface{}{"action": "create"}

	_, err := invokeLambda(context.Background(), client, "test-function", payload)
	if err == nil {
		t.Error("expected error on server failure")
	}
}

// --- CRUD Function Tests ---

func TestLambdaInvokeResourceCreate_Success(t *testing.T) {
	withMockLambdaClient(lambdaSuccessHandler(`{"status": "created"}`), func() {
		d := newLambdaTestResourceData()

		diags := resourceSumologicLambdaInvokeActionCreate(context.Background(), d, nil)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		if d.Id() == "" {
			t.Error("expected ID to be set")
		}
		if !strings.HasPrefix(d.Id(), "my-lambda-") {
			t.Errorf("expected ID to start with 'my-lambda-', got '%s'", d.Id())
		}
		if d.Get("last_lambda_output").(string) == "" {
			t.Error("expected last_lambda_output to be set")
		}
		if d.Get("last_resource_properties").(string) == "" {
			t.Error("expected last_resource_properties to be set")
		}
	})
}

func TestLambdaInvokeResourceCreate_PropertiesStored(t *testing.T) {
	withMockLambdaClient(lambdaSuccessHandler(`{}`), func() {
		d := newLambdaTestResourceData()
		d.Set("aws_resource", "arn:aws:s3:::bucket")
		d.Set("bucket_name", "target-bucket")
		d.Set("account_id", "999999999999")
		d.Set("remove_on_delete_stack", true)

		diags := resourceSumologicLambdaInvokeActionCreate(context.Background(), d, nil)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		var storedProps map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("last_resource_properties").(string)), &storedProps)
		if err != nil {
			t.Fatalf("failed to parse stored resource properties: %v", err)
		}

		if storedProps["AWSResource"] != "arn:aws:s3:::bucket" {
			t.Errorf("expected AWSResource 'arn:aws:s3:::bucket', got '%v'", storedProps["AWSResource"])
		}
		if storedProps["BucketName"] != "target-bucket" {
			t.Errorf("expected BucketName 'target-bucket', got '%v'", storedProps["BucketName"])
		}
		if storedProps["AccountID"] != "999999999999" {
			t.Errorf("expected AccountID '999999999999', got '%v'", storedProps["AccountID"])
		}
		if storedProps["RemoveOnDeleteStack"] != true {
			t.Errorf("expected RemoveOnDeleteStack true, got '%v'", storedProps["RemoveOnDeleteStack"])
		}
	})
}

func TestLambdaInvokeResourceCreate_LambdaError(t *testing.T) {
	withMockLambdaClient(lambdaErrorHandler(), func() {
		d := newLambdaTestResourceData()

		diags := resourceSumologicLambdaInvokeActionCreate(context.Background(), d, nil)
		if !diags.HasError() {
			t.Error("expected error when Lambda invocation fails")
		}
	})
}

func TestLambdaInvokeResourceCreate_ClientError(t *testing.T) {
	original := newLambdaClientFunc
	newLambdaClientFunc = func(ctx context.Context, region string) (*lambda.Client, error) {
		return nil, fmt.Errorf("AWS credentials not configured")
	}
	defer func() { newLambdaClientFunc = original }()

	d := newLambdaTestResourceData()
	diags := resourceSumologicLambdaInvokeActionCreate(context.Background(), d, nil)
	if !diags.HasError() {
		t.Error("expected error when Lambda client cannot be created")
	}
}

func TestLambdaInvokeResourceRead_NoOp(t *testing.T) {
	d := newLambdaTestResourceData()
	d.SetId("test-id")

	diags := resourceSumologicLambdaInvokeActionRead(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if d.Id() != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", d.Id())
	}
}

func TestLambdaInvokeResourceUpdate_Success(t *testing.T) {
	withMockLambdaClient(lambdaSuccessHandler(`{"status": "updated"}`), func() {
		d := newLambdaTestResourceData()
		d.SetId("my-lambda-123")

		oldProps, _ := json.Marshal(map[string]interface{}{
			"AWSResource": "old-resource",
			"BucketName":  "old-bucket",
		})
		d.Set("last_resource_properties", string(oldProps))
		d.Set("aws_resource", "new-resource")
		d.Set("bucket_name", "new-bucket")

		diags := resourceSumologicLambdaInvokeActionUpdate(context.Background(), d, nil)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		if d.Id() != "my-lambda-123" {
			t.Errorf("expected ID preserved as 'my-lambda-123', got '%s'", d.Id())
		}
		if d.Get("last_lambda_output").(string) == "" {
			t.Error("expected last_lambda_output to be set")
		}

		var storedProps map[string]interface{}
		json.Unmarshal([]byte(d.Get("last_resource_properties").(string)), &storedProps)
		if storedProps["AWSResource"] != "new-resource" {
			t.Errorf("expected AWSResource 'new-resource', got '%v'", storedProps["AWSResource"])
		}
	})
}

func TestLambdaInvokeResourceUpdate_LambdaError(t *testing.T) {
	withMockLambdaClient(lambdaErrorHandler(), func() {
		d := newLambdaTestResourceData()
		d.SetId("my-lambda-123")
		d.Set("last_resource_properties", "{}")

		diags := resourceSumologicLambdaInvokeActionUpdate(context.Background(), d, nil)
		if !diags.HasError() {
			t.Error("expected error when Lambda invocation fails")
		}
	})
}

func TestLambdaInvokeResourceDelete_Success(t *testing.T) {
	withMockLambdaClient(lambdaSuccessHandler(`{"status": "deleted"}`), func() {
		d := newLambdaTestResourceData()
		d.SetId("my-lambda-123")

		resourceProps, _ := json.Marshal(map[string]interface{}{
			"AWSResource":         "my-resource",
			"BucketName":          "my-bucket",
			"RemoveOnDeleteStack": true,
		})
		d.Set("last_resource_properties", string(resourceProps))

		diags := resourceSumologicLambdaInvokeActionDelete(context.Background(), d, nil)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		if d.Id() != "" {
			t.Errorf("expected ID to be cleared, got '%s'", d.Id())
		}
	})
}

func TestLambdaInvokeResourceDelete_LambdaError(t *testing.T) {
	withMockLambdaClient(lambdaErrorHandler(), func() {
		d := newLambdaTestResourceData()
		d.SetId("my-lambda-123")
		d.Set("last_resource_properties", "{}")

		diags := resourceSumologicLambdaInvokeActionDelete(context.Background(), d, nil)
		if !diags.HasError() {
			t.Error("expected error when Lambda invocation fails")
		}
	})
}

func TestLambdaInvokeResourceDelete_ClientError(t *testing.T) {
	original := newLambdaClientFunc
	newLambdaClientFunc = func(ctx context.Context, region string) (*lambda.Client, error) {
		return nil, fmt.Errorf("AWS credentials not configured")
	}
	defer func() { newLambdaClientFunc = original }()

	d := newLambdaTestResourceData()
	d.SetId("my-lambda-123")
	diags := resourceSumologicLambdaInvokeActionDelete(context.Background(), d, nil)
	if !diags.HasError() {
		t.Error("expected error when Lambda client cannot be created")
	}
}
