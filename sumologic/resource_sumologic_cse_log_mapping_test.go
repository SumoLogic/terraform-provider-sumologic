package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCELogMapping_create(t *testing.T) {
	SkipCseTest(t)

	var logMapping CSELogMapping
	lmName := "New Log Mapping"
	lmRecordType := "Audit"
	lmEnabled := true
	lmRelatesEntities := true
	lmSkippedValue := "skipped"
	lmProduct := "Web Gateway"
	lmVendor := "McAfee"
	lmLookUp := CSELogMappingLookUp{
		Key:   "tunnel-up",
		Value: "true",
	}

	lmField := CSELogMappingField{
		Name:             "action",
		Value:            "action",
		ValueType:        "constant",
		SkippedValues:    []string{"-"},
		Format:           "JSON",
		CaseInsensitive:  false,
		AlternateValues:  []string{"altValue"},
		TimeZone:         "UTC",
		SplitDelimiter:   ",",
		SplitIndex:       "0",
		FieldJoin:        []string{"and"},
		JoinDelimiter:    "",
		FormatParameters: []string{"param"},
	}

	lmStructuredInputsFields :=
		CSELogMappingStructuredInputField{
			EventIdPattern: "vpn",
			LogFormat:      "JSON",
		}

	resourceName := "sumologic_cse_log_mapping.log_mapping"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSELogMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSELogMappingConfig(lmName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSELogMappingExists(resourceName, &logMapping),
					testCheckLogMappingValues(&logMapping, lmName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccSumologicSCELogMapping_update(t *testing.T) {
	SkipCseTest(t)

	var logMapping CSELogMapping
	lmName := "New Log Mapping"
	lmRecordType := "Audit"
	lmEnabled := true
	lmRelatesEntities := true
	lmSkippedValue := "skipped"
	lmProduct := "Web Gateway"
	lmVendor := "McAfee"
	lmLookUp := CSELogMappingLookUp{
		Key:   "tunnel-up",
		Value: "true",
	}

	lmField := CSELogMappingField{
		Name:             "action",
		Value:            "action",
		ValueType:        "constant",
		SkippedValues:    []string{"-"},
		Format:           "JSON",
		CaseInsensitive:  false,
		AlternateValues:  []string{"altValue"},
		TimeZone:         "UTC",
		SplitDelimiter:   ",",
		SplitIndex:       "0",
		FieldJoin:        []string{"and"},
		JoinDelimiter:    "",
		FormatParameters: []string{"param"},
	}

	lmStructuredInputsFields :=
		CSELogMappingStructuredInputField{
			EventIdPattern: "vpn",
			LogFormat:      "JSON",
		}

	uName := "Changed Name"
	resourceName := "sumologic_cse_log_mapping.log_mapping"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSELogMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSELogMappingConfig(lmName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSELogMappingExists(resourceName, &logMapping),
					testCheckLogMappingValues(&logMapping, lmName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSELogMappingConfig(uName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSELogMappingExists(resourceName, &logMapping),
					testCheckLogMappingValues(&logMapping, uName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValue, lmField, lmLookUp, lmStructuredInputsFields, lmProduct, lmVendor),
				),
			},
		},
	})
}

func testAccCSELogMappingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_log_mapping" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Log Mapping destruction check: CSE Log Mapping ID is not set")
		}

		s, err := client.GetCSELogMapping(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("log mapping still exists")
		}
	}
	return nil
}

func testCreateCSELogMappingConfig(lmName string, lmRecordType string, lmEnabled bool, lmRelatesEntities bool, lmSkippedValues string, lmField CSELogMappingField, lmLookUp CSELogMappingLookUp, lmStructuredInputsFields CSELogMappingStructuredInputField, lmProduct string, lmVendor string) string {

	resource := fmt.Sprintf(`

data "sumologic_cse_log_mapping_vendor_product" "web_gateway" {
  	product = "%s"
	vendor = "%s"
}

resource "sumologic_cse_log_mapping" "log_mapping" {
	name = "%s"
	product_guid = "${data.sumologic_cse_log_mapping_vendor_product.web_gateway.guid}"
	record_type = "%s"
	enabled = "%t"
	relates_entities = "%t"
	skipped_values = ["%s"]
	fields {
			name = "%s"
			value = "%s"
			value_type = "%s"
			skipped_values = ["%s"]
			default_value = "%s"
			format = "%s"
			case_insensitive = "%t"
			alternate_values = ["%s"]
			time_zone = "%s"
			split_delimiter = "%s"
			split_index = "%s"
			field_join = ["%s"]
			join_delimiter = "%s"
			format_parameters = ["%s"]
			lookup {
					key = "%s"
					value = "%s"
			}
		}
	structured_inputs  {
			event_id_pattern = "%s"
			log_format = "%s"
			product = "${data.sumologic_cse_log_mapping_vendor_product.web_gateway.product}"
			vendor = "${data.sumologic_cse_log_mapping_vendor_product.web_gateway.vendor}"
	}
}



`, lmProduct, lmVendor, lmName, lmRecordType, lmEnabled, lmRelatesEntities, lmSkippedValues,
		lmField.Name, lmField.Value, lmField.ValueType, lmField.SkippedValues[0], lmField.DefaultValue, lmField.Format, lmField.CaseInsensitive, lmField.AlternateValues[0], lmField.TimeZone, lmField.SplitDelimiter, lmField.SplitIndex, lmField.FieldJoin[0], lmField.JoinDelimiter, lmField.FormatParameters[0],
		lmLookUp.Key, lmLookUp.Value,
		lmStructuredInputsFields.EventIdPattern, lmStructuredInputsFields.LogFormat)

	return resource
}

func testCheckCSELogMappingExists(n string, logMapping *CSELogMapping) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("log mapping ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		logMappingResp, err := c.GetCSELogMapping(rs.Primary.ID)
		if err != nil {
			return err
		}

		*logMapping = *logMappingResp

		return nil
	}
}

func testCheckLogMappingValues(logMapping *CSELogMapping, lmName string, lmRecordType string, lmEnabled bool, lmRelatesEntities bool, lmSkippedValues string, lmField CSELogMappingField, lmLookUp CSELogMappingLookUp, lmStructuredInputsFields CSELogMappingStructuredInputField, lmProduct string, lmVendor string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logMapping.Name != lmName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", lmName, logMapping.Name)
		}
		if logMapping.RecordType != lmRecordType {
			return fmt.Errorf("bad record type, expected \"%s\", got: %#v", lmRecordType, logMapping.RecordType)
		}
		if *(logMapping.Enabled) != lmEnabled {
			return fmt.Errorf("bad enabled flag, expected \"%t\", got: %#v", lmEnabled, logMapping.Enabled)
		}
		if logMapping.RelatesEntities != lmRelatesEntities {
			return fmt.Errorf("bad relatesEntities flag, expected \"%t\", got: %#v", lmRelatesEntities, logMapping.RelatesEntities)
		}
		if logMapping.SkippedValues[0] != lmSkippedValues {
			return fmt.Errorf("bad skippedValues, expected \"%s\", got: %#v", lmSkippedValues, logMapping.SkippedValues[0])
		}
		if logMapping.Fields[0].Name != lmField.Name ||
			logMapping.Fields[0].Value != lmField.Value ||
			logMapping.Fields[0].ValueType != lmField.ValueType ||
			logMapping.Fields[0].SkippedValues[0] != lmField.SkippedValues[0] ||
			logMapping.Fields[0].DefaultValue != lmField.DefaultValue ||
			logMapping.Fields[0].Format != lmField.Format ||
			logMapping.Fields[0].CaseInsensitive != lmField.CaseInsensitive ||
			logMapping.Fields[0].AlternateValues[0] != lmField.AlternateValues[0] ||
			logMapping.Fields[0].TimeZone != lmField.TimeZone ||
			logMapping.Fields[0].SplitDelimiter != lmField.SplitDelimiter ||
			logMapping.Fields[0].SplitIndex != lmField.SplitIndex ||
			logMapping.Fields[0].FieldJoin[0] != lmField.FieldJoin[0] ||
			logMapping.Fields[0].JoinDelimiter != lmField.JoinDelimiter ||
			logMapping.Fields[0].LookUp[0].Key != lmLookUp.Key || logMapping.Fields[0].LookUp[0].Value != lmLookUp.Value {

			return fmt.Errorf("bad field, expected \"%#v\", got: %#v", lmField, logMapping.Fields[0])
		}

		if logMapping.StructuredInputs[0].Product != lmProduct ||
			logMapping.StructuredInputs[0].Vendor != lmVendor ||
			logMapping.StructuredInputs[0].LogFormat != lmStructuredInputsFields.LogFormat ||
			logMapping.StructuredInputs[0].EventIdPattern != lmStructuredInputsFields.EventIdPattern {

			return fmt.Errorf("bad structured input, expected \"%#v\"", logMapping.StructuredInputs[0])

		}

		return nil

	}
}
