package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSELogMapping(id string) (*CSELogMapping, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/log-mappings/%s", id), false)
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

	responseBody, err := s.Post("sec/v1/log-mappings", request, false)
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
	request := CSELogMappingRequest{
		CSELogMapping: CSELogMapping,
	}

	_, err := s.Put(url, request, false)

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
	Enabled            bool                                `json:"enabled"`
	RelatesEntities    bool                                `json:"relatesEntities"`
	Fields             []CSELogMappingField                `json:"fields"`
	SkippedValues      []string                            `json:"skippedValues"`
	StructuredFields   CSELogMappingStructuredInputField   `json:"structuredFields"`
	StructuredInputs   []CSELogMappingStructuredInputField `json:"structuredInputs"`
	UnstructuredFields CSELogMappingUnstructuredFields     `json:"unstructuredFields"`
}

type CSELogMappingField struct {
	Name             string                `json:"name"`
	Value            string                `json:"value"`
	ValueType        string                `json:"valueType"`
	SkippedValues    []string              `json:"skippedValues"`
	DefaultValue     string                `json:"defaultValue"`
	Format           string                `json:"format"`
	CaseInsensitive  bool                  `json:"caseInsensitive"`
	AlternateValues  []string              `json:"alternateValues"`
	TimeZone         string                `json:"timeZone"`
	SplitDelimiter   string                `json:"splitDelimiter"`
	SplitIndex       string                `json:"splitIndex"`
	FieldJoin        []string              `json:"fieldJoin"`
	JoinDelimiter    string                `json:"joinDelimiter"`
	FormatParameters []string              `json:"formatParameters"`
	LookUp           []CSELogMappingLookUp `json:"lookUp"`
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
	PatternNames []string `json:"patternNames"`
}
