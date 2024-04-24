package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func GetCmfFgpPermStmtSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subject_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"subject_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"role", "org"}, false),
		},
		"permissions": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice(
					[]string{"Create", "Read", "Update", "Delete", "Manage"}, false),
			},
			Required: true,
		},
	}
}

func GetCmfFgpObjPermSetSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: GetCmfFgpPermStmtSchema(),
		},
		// NOTE(2022-05-04): ValidateFunc is not yet supported on lists or sets
		Optional: true,
	}
}

func ResourceToCmfFgpPermStmts(d *schema.ResourceData, targetId string) ([]CmfFgpPermStatement, error) {
	permStmtResourceList := d.Get("obj_permission").(*schema.Set).List()

	var result []CmfFgpPermStatement
	for i := range permStmtResourceList {
		permStmtMap := permStmtResourceList[i].(map[string]interface{})

		result = append(result,
			CmfFgpPermStatement{
				SubjectId:   permStmtMap["subject_id"].(string),
				SubjectType: permStmtMap["subject_type"].(string),
				TargetId:    targetId,
				Permissions: permissionSetToStringArray(permStmtMap["permissions"]),
			})
	}

	// The following logic should be moved to a ValidateFunc, once ValidateFunc is supported for TypeSet
	tfResourcePermIDSet := make(map[CmfFgpIdentifier]bool)
	for i := range result {
		fgpIdentifier := result[i].ToCmfFgpIdentifier()
		if tfResourcePermIDSet[fgpIdentifier] {
			return nil, fmt.Errorf("duplicated FGP Stmts involving %+v", fgpIdentifier)
		} else {
			tfResourcePermIDSet[fgpIdentifier] = true
		}
	}

	return result, nil
}

func CmfFgpPermStmtsSetToResource(d *schema.ResourceData, permStmts []CmfFgpPermStatement) {

	var permStmtResources []map[string]interface{}
	for i := range permStmts {
		permStmt := permStmts[i]
		permStmtResource := make(map[string]interface{})
		permStmtResource["subject_id"] = permStmt.SubjectId
		permStmtResource["subject_type"] = permStmt.SubjectType
		permStmtResource["permissions"] = permStmt.Permissions
		permStmtResources = append(permStmtResources, permStmtResource)
	}
	d.Set("obj_permission", permStmtResources)
}

func ReconcileFgpPermStmtsWithEmptyPerms(tfResourcePermStmts []CmfFgpPermStatement,
	readPermStmts []CmfFgpPermStatement) []CmfFgpPermStatement {
	tfResourcePermIDSet := make(map[CmfFgpIdentifier]bool)
	for i := range tfResourcePermStmts {
		tfResourcePermIDSet[tfResourcePermStmts[i].ToCmfFgpIdentifier()] = true
	}
	var emptyPermStmts []CmfFgpPermStatement
	for i := range readPermStmts {
		readPermStmt := readPermStmts[i]
		if !tfResourcePermIDSet[readPermStmt.ToCmfFgpIdentifier()] {
			emptyPermStmts = append(emptyPermStmts,
				CmfFgpPermStatement{
					SubjectId:   readPermStmt.SubjectId,
					SubjectType: readPermStmt.SubjectType,
					TargetId:    readPermStmt.TargetId,
					Permissions: make([]string, 0),
					// using empty "Permissions" array to effectively revoke FGP for this Subject
				})
		}
	}
	if len(emptyPermStmts) == 0 {
		return tfResourcePermStmts
	} else {
		return append(tfResourcePermStmts, emptyPermStmts...)
	}

}

func permissionSetToStringArray(permSetIntf interface{}) []string {
	permIntfList := permSetIntf.(*schema.Set).List()
	result := make([]string, len(permIntfList))
	for i := range permIntfList {
		result[i] = permIntfList[i].(string)
	}
	return result
}
