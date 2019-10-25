package sumologic

var FieldsMap = map[string]map[string]string{
	"User": map[string]string{
		"firstName":        "Test",
		"lastName":         "User",
		"email":            "testterraform@demo.com",
		"roleIds":          "000000000000076A",
		"isActive":         "true",
		"updatedFirstName": "TestUpdated",
		"updatedLastName":  "UserUpdated",
		"updatedEmail":     "testterraform@demo.com",
		"updatedRoleIds":   "000000000000076A",
		"updatedIsActive":  "false",
	},

	"Role": map[string]string{
		"name":                   "TestTestTestTerraformRole",
		"description":            "test terraform role",
		"filterPredicate":        "_sourceCategory=org-service",
		"capabilities":           "manageContent",
		"updatedName":            "TestTestTestTerraformRoleUpdated",
		"updatedDescription":     "test terraform role updated",
		"updatedFilterPredicate": "_sourceCategory=bill",
		"updatedCapabilities":    "manageContent",
	},

	"FieldExtractionRule": map[string]string{
		"name":                   "TestTestTestTerraformFER",
		"scope":                  "_sourceHost=127.0.0.1",
		"parseExpression":        "csv _raw extract 1 as f1",
		"enabled":                "true",
		"updatedName":            "TestTestTestTerraformFERUpdated",
		"updatedScope":           "_sourceHost=127.0.0.1",
		"updatedParseExpression": "csv _raw extract 1 as f1",
		"updatedEnabled":         "true",
	},

	"IngestBudget": map[string]string{
		"name":                  "TestTestTerraformIngestBudget",
		"fieldValue":            "testTerraformAcceptanceTest",
		"capacityBytes":         "1000",
		"timezone":              "America/Los_Angeles",
		"resetTime":             "14:10",
		"description":           "string",
		"action":                "keepCollecting",
		"auditThreshold":        "85",
		"updatedName":           "TestTestTerraformIngestBudgetUpdated",
		"updatedFieldValue":     "testTerraformAcceptanceTest",
		"updatedCapacityBytes":  "1000",
		"updatedTimezone":       "America/Los_Angeles",
		"updatedResetTime":      "14:10",
		"updatedDescription":    "string updated",
		"updatedAction":         "keepCollecting",
		"updatedAuditThreshold": "90",
	},
}
