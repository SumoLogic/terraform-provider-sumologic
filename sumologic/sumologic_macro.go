package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateMacro(macroReq Macro) (*Macro, error) {
	responseBody, err := s.Post("v2/macros", macroReq)
	if err != nil {
		return nil, err
	}

	var macro Macro
	err = json.Unmarshal(responseBody, &macro)
	if err != nil {
		return nil, err
	}
	log.Printf("[CreateMacro] response: %+v\n", macro)
	return &macro, nil

}

func (s *Client) GetMacro(id string) (*Macro, error) {
	url := fmt.Sprintf("v2/macros/%s", id)
	data, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var macro Macro
	err = json.Unmarshal(data, &macro)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetMacro] response: %+v\n", macro)
	return &macro, nil
}

func (s *Client) DeleteMacro(id string) error {
	url := fmt.Sprintf("v2/macros/%s", id)
	_, err := s.Delete(url)
	return err
}

func (s *Client) UpdateMacro(macro Macro) error {
	url := fmt.Sprintf("v2/macros/%s", macro.ID)
	_, err := s.Put(url, macro)
	return err
}

// ---------- TYPES ----------
type Macro struct {
	ID                  string               `json:"id,omitempty"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	Definition          string               `json:"definition"`
	Enabled             bool                 `json:"enabled"`
	Arguments           []Argument           `json:"arguments"`
	ArgumentValidations []Argumentvalidation `json:"argumentValidations"`
}

type Argument struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Argumentvalidation struct {
	EvalExpression string `json:"evalExpression"`
	ErrorMessage   string `json:"errorMessage"`
}

// ---------- END ----------
