package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSELogMapping(id string) (*CSELogMapping, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/log-mappings/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSELogMappingResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSELogMapping, nil
}

func (s *Client) DeleteCSELogMapping(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/log-mappings/%s", id))

	return err
}

func (s *Client) CreateCSELogMapping(CSELogMapping CSELogMapping) (string, error) {

	request := CSELogMappingRequest{
		CSELogMapping: CSELogMapping,
	}

	var response CSELogMappingResponse

	responseBody, err := s.Post("sec/v1/log-mappings", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSELogMapping.ID, nil
}

func (s *Client) UpdateCSELogMapping(CSELogMapping CSELogMapping) error {
	url := fmt.Sprintf("sec/v1/log-mappings/%s", CSELogMapping.ID)

	CSELogMapping.ID = ""
	CSELogMapping.Enabled = nil
	request := CSELogMappingRequest{
		CSELogMapping: CSELogMapping,
	}

	_, err := s.Put(url, request)

	return err
}

type CSELogMappingRequest struct {
	CSELogMapping CSELogMapping `json:"fields"`
}

type CSELogMappingResponse struct {
	CSELogMapping CSELogMapping `json:"data"`
}

type CSELogMapping struct {
	ID                 string                              `json:"id,omitempty"`
	Name               string                              `json:"name"`
	ParentId           string                              `json:"parentId,omitempty"`
	ProductGuid        string                              `json:"productGuid,omitempty"`
	RecordType         string                              `json:"recordType"`
	RelatesEntities    bool                                `json:"relatesEntities"`
	Fields             []CSELogMappingField                `json:"fields"`
	SkippedValues      []string                            `json:"skippedValues,omitempty"`
	StructuredInputs   []CSELogMappingStructuredInputField `json:"structuredInputs,omitempty"`
	UnstructuredFields *CSELogMappingUnstructuredFields    `json:"unstructuredFields,omitempty"`
	Enabled            *bool                               `json:"enabled,omitempty"`
}

type CSELogMappingField struct {
	Name             string                 `json:"name"`
	Value            string                 `json:"value,omitempty"`
	ValueType        string                 `json:"valueType,omitempty"`
	SkippedValues    []string               `json:"skippedValues,omitempty"`
	DefaultValue     string                 `json:"defaultValue,omitempty"`
	Format           string                 `json:"format,omitempty"`
	CaseInsensitive  bool                   `json:"caseInsensitive,omitempty"`
	AlternateValues  []string               `json:"alternateValues,omitempty"`
	TimeZone         string                 `json:"timeZone,omitempty"`
	SplitDelimiter   string                 `json:"splitDelimiter,omitempty"`
	SplitIndex       string                 `json:"splitIndex,omitempty"`
	FieldJoin        []string               `json:"fieldJoin,omitempty"`
	JoinDelimiter    string                 `json:"joinDelimiter,omitempty"`
	FormatParameters []string               `json:"formatParameters,omitempty"`
	LookUp           *[]CSELogMappingLookUp `json:"lookup,omitempty"`
}

type CSELogMappingLookUp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CSELogMappingStructuredInputField struct {
	EventIdPattern string `json:"eventIdPattern"`
	LogFormat      string `json:"logFormat"`
	Product        string `json:"product"`
	Vendor         string `json:"vendor"`
}

type CSELogMappingUnstructuredFields struct {
	PatternNames []string `json:"patternNames,omitempty"`
}
