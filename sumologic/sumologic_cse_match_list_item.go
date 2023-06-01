package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEMatchListItem(id string) (*CSEMatchListItemGet, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/match-list-items/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEMatchListItemResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEMatchListItemGet, nil
}

func (s *Client) GetCSEMatchListItemsInMatchList(MatchListId string) (*CSEMatchListItemsInMatchListGet, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/match-list-items?listIds=%s", MatchListId))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEMatchListItemsInMatchListResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEMatchListItemsGetData, nil
}

func (s *Client) DeleteCSEMatchListItem(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/match-list-items/%s", id))

	return err
}

func (s *Client) SendCreateCSEMatchListItemsRequest(CSEMatchListItemPost []CSEMatchListItemPost, MatchListID string) error {
	request := CSEMatchListItemRequestPost{
		CSEMatchListItemPost: CSEMatchListItemPost,
	}

	var response CSEMatchListItemResponse

	responseBody, err := s.Post(fmt.Sprintf("sec/v1/match-lists/%s/items", MatchListID), request)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return err
	}

	return nil
}

func (s *Client) CreateCSEMatchListItems(CSEMatchListItemPost []CSEMatchListItemPost, MatchListID string) error {
	var start = 0
	var end = 1000

	//If there are more than 1000 items, send requests in batches of 1000 due to the API's maximum item limit allowed per request.
	for end >= len(CSEMatchListItemPost) {
		err := s.SendCreateCSEMatchListItemsRequest(CSEMatchListItemPost[start:end], MatchListID)
		if err != nil {
			return err
		}
		start += 1000
		end += 1000
	}

	return s.SendCreateCSEMatchListItemsRequest(CSEMatchListItemPost[start:], MatchListID)

}

func (s *Client) UpdateCSEMatchListItem(CSEMatchListItemPost CSEMatchListItemPost) error {
	url := fmt.Sprintf("sec/v1/match-list-items/%s", CSEMatchListItemPost.ID)

	request := CSEMatchListItemRequestUpdate{
		CSEMatchListItemUpdate{
			Active:      CSEMatchListItemPost.Active,
			Expiration:  CSEMatchListItemPost.Expiration,
			Description: CSEMatchListItemPost.Description,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSEMatchListItemRequestPost struct {
	CSEMatchListItemPost []CSEMatchListItemPost `json:"items"`
}

type CSEMatchListItemRequestUpdate struct {
	CSEMatchListItemUpdate CSEMatchListItemUpdate `json:"fields"`
}

type CSEMatchListItemResponse struct {
	CSEMatchListItemGet CSEMatchListItemGet `json:"data"`
}

type CSEMatchListItemsInMatchListResponse struct {
	CSEMatchListItemsGetData CSEMatchListItemsInMatchListGet `json:"data"`
}

type CSEMatchListItemsInMatchListGet struct {
	CSEMatchListItemsGetObjects []CSEMatchListItemGet `json:"objects"`
	Total                       int                   `json:"total"`
}

type CSEMatchListItemPost struct {
	ID          string `json:"id,omitempty"`
	Active      bool   `json:"active,omitempty"`
	Description string `json:"description,omitempty"`
	Expiration  string `json:"expiration,omitempty"`
	Value       string `json:"value,omitempty"`
}

type CSEMatchListItemGet struct {
	ID         string               `json:"id,omitempty"`
	Active     bool                 `json:"active,omitempty"`
	Expiration string               `json:"expiration,omitempty"`
	Meta       CSEMatchListItemMeta `json:"meta,omitempty"`
	Value      string               `json:"value,omitempty"`
}

type CSEMatchListItemsGetRequest struct {
	ID         string               `json:"id,omitempty"`
	Active     bool                 `json:"active,omitempty"`
	Expiration string               `json:"expiration,omitempty"`
	Meta       CSEMatchListItemMeta `json:"meta,omitempty"`
	Value      string               `json:"value,omitempty"`
}

type CSEMatchListItemUpdate struct {
	Active      bool   `json:"active,omitempty"`
	Expiration  string `json:"expiration,omitempty"`
	Description string `json:"description,omitempty"`
}

type CSEMatchListItemMeta struct {
	Created     CSEMatchListItemMetaCreatedUpdated `json:"created,omitempty"`
	Description string                             `json:"description,omitempty"`
	Updated     CSEMatchListItemMetaCreatedUpdated `json:"updated,omitempty"`
}

type CSEMatchListItemMetaCreatedUpdated struct {
	Username string `json:"username,omitempty"`
	When     string `json:"when,omitempty"`
}
