package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSENetworkBlock(id string) (*CSENetworkBlock, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/network-blocks/%s", id), false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSENetworkBlockResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSENetworkBlock, nil
}

func (s *Client) DeleteCSENetworkBlock(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/network-blocks/%s", id))

	return err
}

func (s *Client) CreateCSENetworkBlock(cseNetworkBlock CSENetworkBlock) (string, error) {

	request := CSENetworkBlockRequest{
		CSENetworkBlock: cseNetworkBlock,
	}

	var response CSENetworkBlockResponse

	responseBody, err := s.Post("sec/v1/network-blocks", request, false)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSENetworkBlock.ID, nil
}

func (s *Client) UpdateCSENetworkBlock(cseNetworkBlock CSENetworkBlock) error {
	url := fmt.Sprintf("sec/v1/network-blocks/%s", cseNetworkBlock.ID)
	cseNetworkBlock.ID = ""
	request := CSENetworkBlockRequest{
		CSENetworkBlock: cseNetworkBlock,
	}

	_, err := s.Put(url, request, false)

	return err
}

type CSENetworkBlockRequest struct {
	CSENetworkBlock CSENetworkBlock `json:"fields"`
}

type CSENetworkBlockResponse struct {
	CSENetworkBlock CSENetworkBlock `json:"data"`
}

type CSENetworkBlock struct {
	ID                string `json:"id,omitempty"`
	AddressBlock      string `json:"addressBlock"`
	Label             string `json:"label"`
	Internal          bool   `json:"internal"`
	SuppressesSignals bool   `json:"suppressesSignals"`
}
