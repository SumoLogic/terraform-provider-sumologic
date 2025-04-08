package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

func (s *Client) GetCmfFgp(targetType string, targetId string) (*CmfFgpResponse, error) {

	// e.g. "v1/monitors/0000000000000003/permissions"
	url := fmt.Sprintf("v1/%s/%s/permissions", targetType, targetId)
	data, err := s.GetWithErrOpt(url, true)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var cmfFgpResponse CmfFgpResponse
	err = json.Unmarshal(data, &cmfFgpResponse)
	if err != nil {
		return nil, err
	}

	// Filter subjectType "user"
	cmfFgpResponse.PermissionStatements = cmfFgpFilter(cmfFgpResponse.PermissionStatements,
		func(perm *CmfFgpPermStatement) bool {
			return perm.SubjectType != "user"
		})
	return &cmfFgpResponse, nil
}

func (s *Client) SetCmfFgp(targetType string, cmfFgpRequest CmfFgpRequest) (*CmfFgpResponse, error) {

	if len(cmfFgpRequest.PermissionStatements) == 0 {
		log.Printf("[INFO] SetCmfFgp does not contain any PermissionStatements. Hence, No-Op")
		return nil, nil
	}

	// e.g. "v1/monitors/permissions/set"
	url := fmt.Sprintf("v1/%s/permissions/set", targetType)
	data, err := s.Put(url, cmfFgpRequest)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var cmfFgpResponse CmfFgpResponse
	err = json.Unmarshal(data, &cmfFgpResponse)
	if err != nil {
		return nil, err
	}

	// Filter subjectType "user"
	cmfFgpResponse.PermissionStatements = cmfFgpFilter(cmfFgpResponse.PermissionStatements,
		func(perm *CmfFgpPermStatement) bool {
			return perm.SubjectType != "user"
		})

	return &cmfFgpResponse, nil
}

type CmfFgpResponse struct {
	PermissionStatements []CmfFgpPermStatement `json:"permissionStatements"`
}

type CmfFgpRequest struct {
	PermissionStatements []CmfFgpPermStatement `json:"permissionStatementDefinitions"`
}

type CmfFgpIdentifier struct {
	SubjectId   string `json:"subjectId"`
	SubjectType string `json:"subjectType"`
	TargetId    string `json:"targetId"`
}

// CmfFgpPermStatement is used to handle JSON Marshal and Unmarshal
// for both FGP Request and FGP Response
// Under FGP Response, those "createdAt", "createdBy", "modifiedAt", "modifiedBy" fields
// would be ignored.
type CmfFgpPermStatement struct {
	SubjectId   string   `json:"subjectId"`
	SubjectType string   `json:"subjectType"`
	TargetId    string   `json:"targetId"`
	Permissions []string `json:"permissions"`
}

func CmfFgpPermStmtSetEqual(permStmts01 []CmfFgpPermStatement, permStmts02 []CmfFgpPermStatement) bool {
	permStmtMap01 := cmfFgpPermStmtsToPermStmtMap(permStmts01)
	permStmtMap02 := cmfFgpPermStmtsToPermStmtMap(permStmts02)
	if len(permStmtMap01) != len(permStmtMap02) {
		return false
	}
	for k, v01 := range permStmtMap01 {
		v02, ok := permStmtMap02[k]
		if ok {
			if !cmfFgpStringSetEqual(v01.Permissions, v02.Permissions) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (permStmt *CmfFgpPermStatement) ToCmfFgpIdentifier() CmfFgpIdentifier {
	return CmfFgpIdentifier{
		SubjectId:   permStmt.SubjectId,
		SubjectType: permStmt.SubjectType,
		TargetId:    permStmt.TargetId,
	}
}

func cmfFgpFilter(arr []CmfFgpPermStatement, cond func(*CmfFgpPermStatement) bool) []CmfFgpPermStatement {
	result := []CmfFgpPermStatement{}
	for i := range arr {
		if cond(&(arr[i])) {
			result = append(result, arr[i])
		}
	}
	return result
}

func cmfFgpStringSetEqual(strs01 []string, strs02 []string) bool {
	strmap01 := cmfFgpToStringMap(strs01)
	strmap02 := cmfFgpToStringMap(strs02)

	if len(strmap01) != len(strmap02) {
		return false
	}
	for k := range strmap01 {
		if !strmap02[k] {
			return false
		}
	}
	return true
}

func cmfFgpPermStmtsToPermStmtMap(permStmts []CmfFgpPermStatement) map[CmfFgpIdentifier]CmfFgpPermStatement {
	stmtMap := make(map[CmfFgpIdentifier]CmfFgpPermStatement)
	for _, v := range permStmts {
		stmtMap[v.ToCmfFgpIdentifier()] = v
	}
	return stmtMap
}

func cmfFgpToStringMap(strs []string) map[string]bool {
	strmap := make(map[string]bool)
	for _, v := range strs {
		strmap[v] = true
	}
	return strmap
}

// Under FGP Request:
/*
	{
		"targetId": "0000000000121E9E",
		"subjectType": "org",
		"subjectId": "0000000006B11645",
		"permissions": ["Read"]
	}
*/
// Under FGP Response:
/*
	{
		"permissions": ["Read"],
		"subjectType": "org",
		"subjectId": "0000000006B11645",
		"targetId": "0000000000121E9E",
		"createdAt": "2022-04-27T17:33:53.194Z",
		"createdBy": "0000000006D5CC14",
		"modifiedAt": "2022-04-27T17:33:53.194Z",
		"modifiedBy": "0000000006D5CC14"
	}
*/
