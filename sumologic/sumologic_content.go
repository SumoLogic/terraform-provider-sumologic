package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//READ
func (s *Client) GetContent(id string, timeout time.Duration) (*Content, error) {
	log.Println("####Begin GetContent####")

	url := fmt.Sprintf("v2/content/%s/export", id)
	log.Printf("Content export url: %s", url)

	//Begin the content export job
	rawJID, err := s.Post(url, nil, false)

	//If there was an error, exit here and return it
	if err != nil {
		if strings.Contains(err.Error(), "Content with the given ID does not exist.") {
			return nil, nil
		}
		return nil, err
	}

	//Parse the jobId from the response
	var jid JobId
	err = json.Unmarshal(rawJID, &jid)

	//Exit here if there was an error parsing the json
	if err != nil {
		return nil, err
	}
	log.Printf("JobId: %s", jid.ID)

	url = fmt.Sprintf("v2/content/%s/export/%s/status", id, jid.ID)
	log.Printf("Content export job status url: %s", url)

	//Ensure the job has completed before proceeding
	log.Printf("Job Id: %s", id)
	err = waitForJob(url, timeout, s)
	if err != nil {
		return nil, err
	}

	//Request the results of the job
	var content Content

	url = fmt.Sprintf("v2/content/%s/export/%s/result", id, jid.ID)
	log.Printf("Content export job results url: %s", url)

	rawContent, _, err := s.Get(url, false)

	//Exit here if there was an error during the request
	if err != nil {
		return nil, err
	}

	//Parse the export job results and populate the Content struct
	err = json.Unmarshal(rawContent, &content)

	//Exit here if there was an error parsing the json
	if err != nil {
		return nil, err
	}

	log.Println("Setting content.Config to export results...")
	log.Println(string(rawContent))

	//set the content.Config to be the raw results of the export
	//This should be a complete, correctly formatted json representation of the content
	content.Config = string(rawContent)

	log.Println("####End GetContent####")
	return &content, nil
}

//DELETE
func (s *Client) DeleteContent(id string, timeout time.Duration) error {
	log.Println("####Begin DeleteContent####")

	log.Printf("Deleting Content Id: %s", id)

	url := fmt.Sprintf("v2/content/%s/delete", id)
	log.Printf("Content delete url: %s", url)

	//start the deletion job
	rawJID, err := s.Delete(url)

	if err != nil {
		return err
	}

	var jid JobId

	//Parse the response for the JobId
	err = json.Unmarshal(rawJID, &jid)

	//Exit here if there was an error parsing the json
	if err != nil {
		return err
	}

	url = fmt.Sprintf("v2/content/%s/delete/%s/status", id, jid.ID)
	log.Printf("Content delete job status url: %s", url)

	waitForJob(url, timeout, s)

	log.Println("####End DeleteContent####")
	return err
}

//CREATE or UPDATE
func (s *Client) CreateOrUpdateContent(content Content, timeout time.Duration, overwrite bool) (string, error) {
	log.Println("####Begin CreateOrUpdateContent####")

	url := fmt.Sprintf("v2/content/folders/%s/import?overwrite=%s", content.ParentId, strconv.FormatBool(overwrite))
	log.Printf("Create content url: %s", url)

	//Initiate content creation job
	jobResponse, err := s.PostRawPayload(url, content.Config)

	//Exit if there was an error during the request
	if err != nil {
		return "", err
	}

	//Parse JobId
	var jid JobId
	err = json.Unmarshal(jobResponse, &jid)

	log.Printf("Create Content Job Id: %s", jid.ID)

	//Catch parsing errors
	if err != nil {
		return "", err
	}

	url = fmt.Sprintf("v2/content/folders/%s/import/%s/status", content.ParentId, jid.ID)
	log.Printf("Create content job status url: %s", url)

	waitForJob(url, timeout, s)

	log.Println("####Begin Folder Read####")
	log.Printf("Looking up folder with ID: %s", content.ParentId)

	//build the url
	url = fmt.Sprintf("v2/content/folders/%s", content.ParentId)
	log.Printf("Find folder url: %s", url)

	//Request folder content
	rawParent, _, err := s.Get(url, false)

	if err != nil {
		return "", nil
	}

	var parentContent Content

	//unmarshal
	err = json.Unmarshal(rawParent, &parentContent)
	if err != nil {
		return "", err
	}

	log.Printf("Parent Name: %s", parentContent.Name)
	log.Printf("Child Name: %s", content.Name)
	log.Println("####End Folder Read####")

	log.Println("Begin Searching for Child's name within Parent's children...")

	//Search all of the parent folder's children for a matching content name
	for i := 0; i < len(parentContent.Children); i++ {
		//if the object was found, return the ID
		log.Printf("Found - Childname: %s", parentContent.Children[i].Name)

		// Names must be unique within a folder, so it is okay to match by name here
		if parentContent.Children[i].Name == content.Name {
			log.Println("MATCH")
			log.Printf("ChildId: %s", parentContent.Children[i].ID)
			return parentContent.Children[i].ID, nil
		}

		log.Println("No match")
	}

	log.Println("Content not found, are you searching the right parent?")
	//object wasn't found, should make an error here
	return "", nil
}

func waitForJob(url string, timeout time.Duration, s *Client) error {
	conf := &resource.StateChangeConf{
		Pending: []string{
			"InProgress",
		},
		Target: []string{
			"Success",
		},
		Refresh: func() (interface{}, string, error) {
			log.Println("====Start Job Status Check====")

			var status Status
			b, _, err := s.Get(url, false)
			if err != nil {
				return nil, "", err
			}

			err = json.Unmarshal(b, &status)
			if err != nil {
				return nil, "", err
			}

			log.Printf("Job Status: %s", status.Status)
			log.Printf("Job Message: %s", status.StatusMessage)
			log.Println("Job Errors:")
			log.Println(status.Errors)
			log.Println("====End Job Status Check====")

			if status.Status == "failed" {
				return status, status.Status, fmt.Errorf("job failed: %s", status.StatusMessage)
			}

			return status, status.Status, nil
		},
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := conf.WaitForState()
	return err
}
