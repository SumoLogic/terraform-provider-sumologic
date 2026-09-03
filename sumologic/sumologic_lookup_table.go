package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (s *Client) CreateLookupTable(lookupTable LookupTable) (string, error) {
	urlWithoutParams := "v1/lookupTables"

	data, err := s.Post(urlWithoutParams, lookupTable)
	if err != nil {
		return "", err
	}

	var createdLookupTable LookupTable

	err = json.Unmarshal(data, &createdLookupTable)
	if err != nil {
		return "", err
	}

	log.Printf("##DEBUG## created lookuptable: %+v\n\n", createdLookupTable)
	return createdLookupTable.ID, nil

}

func (s *Client) GetLookupTable(id string) (*LookupTable, error) {
	urlWithoutParams := "v1/lookupTables/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var lookupTable LookupTable

	err = json.Unmarshal(data, &lookupTable)
	if err != nil {
		return nil, err
	}

	return &lookupTable, nil

}

func (s *Client) DeleteLookupTable(id string) error {
	urlWithoutParams := "v1/lookupTables/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	log.Printf("deleting lookuptable: %s", id)
	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateLookupTable(lookupTable LookupTable) error {
	urlWithoutParams := "v1/lookupTables/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, lookupTable.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	lookupTable.ID = ""

	_, err := s.Put(urlWithParams, lookupTable)

	return err

}

type LookupTable struct {
	Fields          []LookupTableField `json:"fields,omitempty"`
	SizeLimitAction string             `json:"sizeLimitAction"`
	Name            string             `json:"name"`
	ID              string             `json:"id,omitempty"`
	ParentFolderId  string             `json:"parentFolderId"`
	PrimaryKeys     []string           `json:"primaryKeys,omitempty"`
	Description     string             `json:"description"`
	Ttl             int                `json:"ttl"`
}

type LookupTableField struct {
	FieldName string `json:"fieldName"`
	FieldType string `json:"fieldType"`
}

// LookupJobStatus mirrors the LookupAsyncJobStatus schema returned by
// v1/lookupTables/jobs/{jobId}/status. It intentionally does not reuse the
// generic Status struct (sumologic_client.go): lookup jobs can report
// Pending/PartialSuccess (which Status/waitForJob treat as errors) and carry
// their error/warning detail in arrays rather than a single message.
type LookupJobStatus struct {
	JobId          string             `json:"jobId"`
	Status         string             `json:"status"`
	StatusMessages []string           `json:"statusMessages,omitempty"`
	Errors         []Error            `json:"errors,omitempty"`
	Warnings       []LookupJobWarning `json:"warnings,omitempty"`
}

type LookupJobWarning struct {
	Message string `json:"message"`
	Cause   string `json:"cause,omitempty"`
}

// UploadLookupTableContent replaces the contents of the lookup table with the
// given CSV. Existing rows are removed (merge=false) since Terraform's
// declarative model treats content as the full desired state of the table.
// The returned warning is non-empty when the job completed as PartialSuccess;
// callers should surface it to the user (e.g. as a diag.Diagnostics warning)
// rather than only logging it, since some rows may have been dropped.
func (s *Client) UploadLookupTableContent(id string, content string, timeout time.Duration) (string, error) {
	urlPath := fmt.Sprintf("v1/lookupTables/%s/upload?merge=false", id)

	data, err := s.PostMultipartFile(urlPath, "file", "lookup.csv", []byte(content))
	if err != nil {
		return "", err
	}

	var jobId JobId
	if err := json.Unmarshal(data, &jobId); err != nil {
		return "", err
	}

	return waitForLookupJob(jobId.ID, timeout, s)
}

// TruncateLookupTable deletes all rows from the lookup table. See
// UploadLookupTableContent for the meaning of the returned warning.
func (s *Client) TruncateLookupTable(id string, timeout time.Duration) (string, error) {
	urlPath := fmt.Sprintf("v1/lookupTables/%s/truncate", id)

	data, err := s.Post(urlPath, nil)
	if err != nil {
		return "", err
	}

	var jobId JobId
	if err := json.Unmarshal(data, &jobId); err != nil {
		return "", err
	}

	return waitForLookupJob(jobId.ID, timeout, s)
}

// waitForLookupJob polls until the job reaches a terminal state. It returns a
// non-empty warning string on PartialSuccess so callers can decide how to
// surface it (this package has no dependency on the SDK's diag package).
func waitForLookupJob(jobId string, timeout time.Duration, s *Client) (string, error) {
	url := fmt.Sprintf("v1/lookupTables/jobs/%s/status", jobId)

	conf := &resource.StateChangeConf{
		Pending: []string{
			"Pending",
			"InProgress",
		},
		Target: []string{
			"Success",
			"PartialSuccess",
		},
		Refresh: func() (interface{}, string, error) {
			var status LookupJobStatus
			b, err := s.Get(url)
			if err != nil {
				return nil, "", err
			}

			if err := json.Unmarshal(b, &status); err != nil {
				return nil, "", err
			}

			if status.Status == "Failed" {
				return status, status.Status, fmt.Errorf("lookup table job failed: %s", formatLookupJobErrors(status.Errors))
			}

			return status, status.Status, nil
		},
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	result, err := conf.WaitForState()
	if err != nil {
		return "", err
	}

	if status, ok := result.(LookupJobStatus); ok && status.Status == "PartialSuccess" {
		warning := formatLookupJobWarnings(status.Warnings)
		if warning == "" {
			warning = "job completed with PartialSuccess but returned no warning detail"
		}
		log.Printf("[WARN] lookup table job %s completed with warnings: %s", jobId, warning)
		return warning, nil
	}

	return "", nil
}

func formatLookupJobErrors(errs []Error) string {
	messages := make([]string, len(errs))
	for i, e := range errs {
		messages[i] = e.Message
	}
	return strings.Join(messages, "; ")
}

func formatLookupJobWarnings(warnings []LookupJobWarning) string {
	messages := make([]string, len(warnings))
	for i, w := range warnings {
		messages[i] = w.Message
	}
	return strings.Join(messages, "; ")
}
