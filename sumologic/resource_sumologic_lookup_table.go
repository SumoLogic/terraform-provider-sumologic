package sumologic

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Sumo Logic's documented limits for the lookup table CSV upload API
// (https://api.sumologic.com/docs/#operation/uploadFile). maxContentBytes is
// enforced as a hard error; warnContentBytes is a softer nudge that large,
// slow-changing data belongs in an ingest pipeline rather than Terraform.
const (
	maxContentBytes  = 100 * 1024 * 1024
	warnContentBytes = 10 * 1024 * 1024
)

func resourceSumologicLookupTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSumologicLookupTableCreate,
		ReadContext:   resourceSumologicLookupTableRead,
		UpdateContext: resourceSumologicLookupTableUpdate,
		DeleteContext: resourceSumologicLookupTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: lookupTableCustomizeDiff,

		Schema: map[string]*schema.Schema{

			"primary_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The primary key field names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},

			"parent_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(0, 1000),
				Required:     true,
			},

			"size_limit_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"StopIncomingMessages", "DeleteOldData"}, false),
				Default:      "StopIncomingMessages",
			},

			"fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of fields in the lookup table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"field_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},

						"field_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validation.StringInSlice([]string{"boolean", "int", "long", "double", "string"}, false),
						},
					},
				},
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "CSV content to populate the lookup table with. The first row must be a header matching the field names in \"fields\" (case-insensitive). Uploading replaces all existing rows. This value cannot be read back from the API: importing an existing table will show a diff on the first plan if \"content\" is set, and the state stores only a hash of the content, not the content itself. Removing this argument truncates the table.",
				StateFunc:    hashLookupContent,
				ValidateFunc: validateLookupContent,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceSumologicLookupTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	if d.Id() == "" {
		lookupTable := resourceToLookupTable(d)
		id, err := c.CreateLookupTable(lookupTable)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)
	}

	var diags diag.Diagnostics

	if content := d.Get("content").(string); content != "" {
		warning, err := c.UploadLookupTableContent(d.Id(), content, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			// The table was already created; return the error without clearing
			// the id so Terraform taints the resource instead of orphaning it.
			return diag.FromErr(err)
		}
		if warning != "" {
			diags = append(diags, lookupTablePartialSuccessDiagnostic(warning))
		}
	}

	log.Printf("created lookup: %+v\n", d)
	log.Printf("lookup id: %v\n", d.Id())
	return append(diags, resourceSumologicLookupTableRead(ctx, d, meta)...)
}

func resourceSumologicLookupTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	id := d.Id()
	lookupTable, err := c.GetLookupTable(id)
	log.Printf("##DEBUG## read lookup: %+v\n", lookupTable)
	if err != nil {
		return diag.FromErr(err)
	}

	if lookupTable == nil {
		log.Printf("[WARN] LookupTable not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", lookupTable.Name)
	if err := d.Set("fields", fieldsToList(lookupTable.Fields)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err))
	}
	d.Set("ttl", lookupTable.Ttl)
	d.Set("primary_keys", lookupTable.PrimaryKeys)
	d.Set("parent_folder_id", lookupTable.ParentFolderId)
	d.Set("size_limit_action", lookupTable.SizeLimitAction)
	d.Set("description", lookupTable.Description)

	// "content" is intentionally not set here: there is no API to read a
	// lookup table's rows back out, so any value written would be a guess.
	// The existing state (or config) value is left as-is.

	return nil
}

func resourceSumologicLookupTableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	log.Printf("##DEBUG## resourceSumologicLookupTableDelete: %s", d.Id())
	return diag.FromErr(c.DeleteLookupTable(d.Id()))
}

func resourceSumologicLookupTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	lookupTable := resourceToLookupTable(d)
	if err := c.UpdateLookupTable(lookupTable); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if d.HasChange("content") {
		var warning string
		var err error
		if content := d.Get("content").(string); content != "" {
			warning, err = c.UploadLookupTableContent(d.Id(), content, d.Timeout(schema.TimeoutUpdate))
		} else {
			// Removing the argument truncates the table rather than leaving
			// stale rows that Terraform no longer has any record of.
			warning, err = c.TruncateLookupTable(d.Id(), d.Timeout(schema.TimeoutUpdate))
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if warning != "" {
			diags = append(diags, lookupTablePartialSuccessDiagnostic(warning))
		}
	}

	return append(diags, resourceSumologicLookupTableRead(ctx, d, meta)...)
}

// lookupTablePartialSuccessDiagnostic surfaces a PartialSuccess job's warning
// text as a non-fatal diagnostic, visible in `terraform apply` output. A
// plain `error` return can't do this - either it fails the apply (wrong: the
// upload did partially succeed) or it says nothing (the prior behavior).
func lookupTablePartialSuccessDiagnostic(warning string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Lookup table upload completed with warnings",
		Detail:   warning,
	}
}

func resourceToLookupTable(d *schema.ResourceData) LookupTable {

	fieldsData := d.Get("fields").([]interface{})
	var fields []LookupTableField
	for _, data := range fieldsData {
		fields = append(fields, resourceToLookupTableField([]interface{}{data}))
	}

	primaryKeysData := d.Get("primary_keys").([]interface{})
	var primaryKeys []string
	for _, data := range primaryKeysData {
		primaryKeys = append(primaryKeys, data.(string))
	}

	return LookupTable{
		Name:            d.Get("name").(string),
		ID:              d.Id(),
		Fields:          fields,
		Description:     d.Get("description").(string),
		Ttl:             d.Get("ttl").(int),
		SizeLimitAction: d.Get("size_limit_action").(string),
		PrimaryKeys:     primaryKeys,
		ParentFolderId:  d.Get("parent_folder_id").(string),
	}
}

func resourceToLookupTableField(data interface{}) LookupTableField {

	lookupTableFieldSlice := data.([]interface{})
	lookupTableField := LookupTableField{}
	if len(lookupTableFieldSlice) > 0 {
		lookupTableFieldObj := lookupTableFieldSlice[0].(map[string]interface{})
		lookupTableField.FieldName = lookupTableFieldObj["field_name"].(string)
		lookupTableField.FieldType = lookupTableFieldObj["field_type"].(string)
	}

	return lookupTableField
}

func fieldsToList(lookupTableField []LookupTableField) []map[string]interface{} {
	var s []map[string]interface{}

	for _, t := range lookupTableField {
		mapping := map[string]interface{}{
			"field_name": t.FieldName,
			"field_type": t.FieldType,
		}
		s = append(s, mapping)
	}

	return s
}

// hashLookupContent is the StateFunc for "content": state stores only a
// sha256 of the CSV, not the CSV itself, since the content can be arbitrarily
// large and cannot be derived back from the API. The empty string is left as
// the empty string so a plan clearing "content" reads as "-> \"\"" (truncate)
// rather than a hash of the empty string.
func hashLookupContent(v interface{}) string {
	content := v.(string)
	if content == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

// validateLookupContent runs cheap, single-value checks on "content" at plan
// time. Checks that need sibling attributes (fields, primary_keys) live in
// lookupTableCustomizeDiff instead, since ValidateFunc only sees one value.
// lintignore:V012
func validateLookupContent(v interface{}, k string) (warnings []string, errors []error) {
	content, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q: expected type string", k))
		return warnings, errors
	}

	if size := len(content); size > maxContentBytes {
		errors = append(errors, fmt.Errorf(
			"%q is %d bytes, which exceeds Sumo Logic's %d byte (100MB) lookup table CSV upload limit",
			k, size, maxContentBytes))
	} else if size > warnContentBytes {
		warnings = append(warnings, fmt.Sprintf(
			"%q is %d bytes; consider populating large lookup tables via an ingest pipeline or scheduled search instead of Terraform, since every plan and apply carries this value across the provider boundary",
			k, size))
	}

	if strings.Contains(content, "\r\n") {
		warnings = append(warnings, fmt.Sprintf(
			"%q contains Windows-style line endings (\\r\\n); Sumo Logic's lookup table upload API expects Unix-style newlines (\\n)",
			k))
	}

	return warnings, errors
}

// lookupTableCustomizeDiff validates the CSV structure of "content" against
// the table's "fields" and "primary_keys" at plan time, so a malformed CSV
// fails terraform plan instead of terraform apply.
func lookupTableCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if !d.NewValueKnown("content") {
		// Content built from an interpolation unresolved at plan time; nothing
		// to validate yet without failing the plan spuriously.
		return nil
	}

	content := d.Get("content").(string)
	if content == "" {
		return nil
	}

	if !d.NewValueKnown("fields") {
		return nil
	}

	fieldsData := d.Get("fields").([]interface{})
	if len(fieldsData) == 0 {
		return nil
	}

	var fields []LookupTableField
	for _, data := range fieldsData {
		fields = append(fields, resourceToLookupTableField([]interface{}{data}))
	}

	primaryKeysData := d.Get("primary_keys").([]interface{})
	var primaryKeys []string
	for _, data := range primaryKeysData {
		primaryKeys = append(primaryKeys, data.(string))
	}

	return validateLookupCSV(content, fields, primaryKeys)
}

// validateLookupCSV checks that content parses as CSV and that its header
// matches fields/primaryKeys. It deliberately does not validate cell values
// against field types: that would duplicate server-side coercion rules this
// provider cannot observe, and getting it wrong would reject configs the API
// would have accepted. Pulled out as a pure function so it is testable
// without a *schema.ResourceDiff.
func validateLookupCSV(content string, fields []LookupTableField, primaryKeys []string) error {
	fieldNames := make(map[string]bool, len(fields))
	for _, f := range fields {
		fieldNames[strings.ToLower(f.FieldName)] = true
	}

	reader := csv.NewReader(strings.NewReader(content))

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to parse CSV header of \"content\": %w", err)
	}

	seenColumns := make(map[string]bool, len(header))
	for i, col := range header {
		trimmed := strings.TrimSpace(col)
		if trimmed == "" {
			return fmt.Errorf("CSV header column %d in \"content\" is empty", i+1)
		}
		lower := strings.ToLower(trimmed)
		if !fieldNames[lower] {
			return fmt.Errorf("CSV header column %q in \"content\" does not match any field name in \"fields\"", trimmed)
		}
		seenColumns[lower] = true
	}

	for _, f := range fields {
		if !seenColumns[strings.ToLower(f.FieldName)] {
			return fmt.Errorf("CSV header in \"content\" is missing column %q, which is defined in \"fields\"", f.FieldName)
		}
	}

	for _, pk := range primaryKeys {
		if !seenColumns[strings.ToLower(pk)] {
			return fmt.Errorf("CSV header in \"content\" is missing column %q, which is listed in \"primary_keys\"", pk)
		}
	}

	// Read the remaining rows so ragged rows and unbalanced quotes fail at
	// plan time. FieldsPerRecord defaults to 0, which locks csv.Reader to the
	// header's column count after the first Read, so mismatched rows are
	// rejected automatically.
	for {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to parse CSV row in \"content\": %w", err)
		}
	}

	return nil
}
