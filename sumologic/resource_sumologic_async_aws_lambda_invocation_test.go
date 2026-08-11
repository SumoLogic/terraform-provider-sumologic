package sumologic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func withMockAsyncLambdaClient(handler http.HandlerFunc, fn func()) {
	server := httptest.NewServer(handler)
	defer server.Close()
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
	original := newAsyncLambdaClientFunc
	newAsyncLambdaClientFunc = func(ctx context.Context, region string, profile string) (*lambda.Client, error) {
		return client, nil
	}
	defer func() { newAsyncLambdaClientFunc = original }()
	fn()
}

func asyncLambdaAcceptedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{}`))
	}
}

func asyncLambdaErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
	}
}

func newAsyncLambdaTestResourceData() *schema.ResourceData {
	r := resourceSumologicAsyncAwsLambdaInvocation()
	d := r.TestResourceData()
	d.Set("function_name", "my-async-lambda")
	d.Set("region", "us-east-1")
	d.Set("input", `{}`)
	d.Set("qualifier", "$LATEST")
	return d
}

// --- Schema Tests ---

func TestAsyncLambdaInvocationResourceSchema(t *testing.T) {
	r := resourceSumologicAsyncAwsLambdaInvocation()

	expectedAttrs := map[string]struct {
		typ      schema.ValueType
		required bool
		optional bool
		computed bool
		forceNew bool
	}{
		"function_name": {typ: schema.TypeString, required: true, forceNew: true},
		"region":        {typ: schema.TypeString, required: true, forceNew: true},
		"input":       {typ: schema.TypeString, optional: true, forceNew: true},
		"qualifier":   {typ: schema.TypeString, optional: true, forceNew: true},
		"triggers":    {typ: schema.TypeMap, optional: true, forceNew: true},
		"aws_profile": {typ: schema.TypeString, optional: true, forceNew: true},
		"status_code": {typ: schema.TypeInt, computed: true},
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
		if attr.ForceNew != expected.forceNew {
			t.Errorf("attribute '%s': expected ForceNew=%v, got %v", name, expected.forceNew, attr.ForceNew)
		}
	}
}

// --- CRUD Tests ---

func TestAsyncLambdaInvocationCreate_Success(t *testing.T) {
	withMockAsyncLambdaClient(asyncLambdaAcceptedHandler(), func() {
		d := newAsyncLambdaTestResourceData()

		diags := resourceAsyncLambdaInvocationCreate(context.Background(), d, nil)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		if d.Id() == "" {
			t.Error("expected ID to be set")
		}
		expectedID := fmt.Sprintf("my-async-lambda-%d", 202)
		if d.Id() != expectedID {
			t.Errorf("expected ID '%s', got '%s'", expectedID, d.Id())
		}
		if d.Get("status_code").(int) != 202 {
			t.Errorf("expected status_code 202, got %d", d.Get("status_code").(int))
		}
	})
}

func TestAsyncLambdaInvocationCreate_LambdaError(t *testing.T) {
	withMockAsyncLambdaClient(asyncLambdaErrorHandler(), func() {
		d := newAsyncLambdaTestResourceData()

		diags := resourceAsyncLambdaInvocationCreate(context.Background(), d, nil)
		if !diags.HasError() {
			t.Error("expected error when Lambda invocation fails")
		}
	})
}

func TestAsyncLambdaInvocationCreate_ClientError(t *testing.T) {
	original := newAsyncLambdaClientFunc
	newAsyncLambdaClientFunc = func(ctx context.Context, region string, profile string) (*lambda.Client, error) {
		return nil, fmt.Errorf("AWS credentials not configured")
	}
	defer func() { newAsyncLambdaClientFunc = original }()

	d := newAsyncLambdaTestResourceData()
	diags := resourceAsyncLambdaInvocationCreate(context.Background(), d, nil)
	if !diags.HasError() {
		t.Error("expected error when Lambda client cannot be created")
	}
}

func TestAsyncLambdaInvocationRead_NoOp(t *testing.T) {
	d := newAsyncLambdaTestResourceData()
	d.SetId("my-async-lambda-202")

	diags := resourceAsyncLambdaInvocationRead(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if d.Id() != "my-async-lambda-202" {
		t.Errorf("expected ID preserved, got '%s'", d.Id())
	}
}

func TestAsyncLambdaInvocationDelete_NoOp(t *testing.T) {
	d := newAsyncLambdaTestResourceData()
	d.SetId("my-async-lambda-202")

	diags := resourceAsyncLambdaInvocationDelete(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
}
