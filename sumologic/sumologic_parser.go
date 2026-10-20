package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const parsersLibraryBaseApiUrl = "v1/parsers"
const parsersLibraryParserType = "ParsersLibraryParser"

// ---------- ENDPOINTS ----------

// GetParsersLibraryRootID returns the identifier of the root folder of the Parsers library.
func (s *Client) GetParsersLibraryRootID() (string, error) {
	data, err := s.Get(parsersLibraryBaseApiUrl + "/root")
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", fmt.Errorf("parsers library root folder not found")
	}

	var root struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return "", err
	}
	return root.ID, nil
}

func (s *Client) CreateParsersLibraryParser(parser ParsersLibraryParser) (string, error) {
	urlWithParams := fmt.Sprintf("%s?parentId=%s", parsersLibraryBaseApiUrl, url.QueryEscape(parser.ParentID))

	data, err := s.Post(urlWithParams, parser.toCreateRequest())
	if err != nil {
		return "", err
	}
	log.Printf("created parser response: %s", string(data))

	var created ParsersLibraryParser
	if err := json.Unmarshal(data, &created); err != nil {
		return "", err
	}
	return created.ID, nil
}

func (s *Client) GetParsersLibraryParser(id string) (*ParsersLibraryParser, error) {
	data, err := s.Get(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var parser ParsersLibraryParser
	if err := json.Unmarshal(data, &parser); err != nil {
		return nil, err
	}
	return &parser, nil
}

// UpdateParsersLibraryParser updates a parser. parser.Version must match the latest
// version of the parser in the library or the API rejects the update.
func (s *Client) UpdateParsersLibraryParser(parser ParsersLibraryParser) error {
	_, err := s.Put(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, parser.ID), parser.toUpdateRequest())
	return err
}

func (s *Client) MoveParsersLibraryParser(id, parentID string) error {
	urlWithParams := fmt.Sprintf("%s/%s/move?parentId=%s", parsersLibraryBaseApiUrl, id, url.QueryEscape(parentID))
	_, err := s.Post(urlWithParams, nil)
	return err
}

func (s *Client) DeleteParsersLibraryParser(id string) error {
	_, err := s.Delete(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, id))
	return err
}

// ---------- TYPES ----------

type ParsersLibraryParser struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Version        int64  `json:"version"`
	CreatedAt      string `json:"createdAt,omitempty"`
	CreatedBy      string `json:"createdBy,omitempty"`
	ModifiedAt     string `json:"modifiedAt,omitempty"`
	ModifiedBy     string `json:"modifiedBy,omitempty"`
	ParentID       string `json:"parentId,omitempty"`
	ContentType    string `json:"contentType,omitempty"`
	Type           string `json:"type"`
	IsLocked       bool   `json:"isLocked"`
	IsSystem       bool   `json:"isSystem"`
	IsMutable      bool   `json:"isMutable"`
	Stanzas        string `json:"stanzas"`
	ModelPath      string `json:"modelPath,omitempty"`
	SourcetypePath string `json:"sourcetypePath,omitempty"`
	Families       string `json:"families,omitempty"`
	IsPartial      bool   `json:"isPartial"`
	LocalStanzas   string `json:"localStanzas,omitempty"`
}

// parsersLibraryParserCreateRequest is the body of POST /v1/parsers (ParsersLibraryParser).
type parsersLibraryParserCreateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	Stanzas        string `json:"stanzas"`
	ModelPath      string `json:"modelPath,omitempty"`
	SourcetypePath string `json:"sourcetypePath,omitempty"`
	Families       string `json:"families,omitempty"`
	IsPartial      bool   `json:"isPartial"`
	LocalStanzas   string `json:"localStanzas,omitempty"`
}

// parsersLibraryParserUpdateRequest is the body of PUT /v1/parsers/{id} (ParsersLibraryParserUpdate).
type parsersLibraryParserUpdateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Version        int64  `json:"version"`
	Type           string `json:"type"`
	Stanzas        string `json:"stanzas"`
	ModelPath      string `json:"modelPath,omitempty"`
	SourcetypePath string `json:"sourcetypePath,omitempty"`
	Families       string `json:"families,omitempty"`
	IsPartial      bool   `json:"isPartial"`
	LocalStanzas   string `json:"localStanzas,omitempty"`
}

func (p ParsersLibraryParser) toCreateRequest() parsersLibraryParserCreateRequest {
	return parsersLibraryParserCreateRequest{
		Name:           p.Name,
		Description:    p.Description,
		Type:           parsersLibraryParserType,
		Stanzas:        p.Stanzas,
		ModelPath:      p.ModelPath,
		SourcetypePath: p.SourcetypePath,
		Families:       p.Families,
		IsPartial:      p.IsPartial,
		LocalStanzas:   p.LocalStanzas,
	}
}

func (p ParsersLibraryParser) toUpdateRequest() parsersLibraryParserUpdateRequest {
	return parsersLibraryParserUpdateRequest{
		Name:           p.Name,
		Description:    p.Description,
		Version:        p.Version,
		Type:           parsersLibraryParserType,
		Stanzas:        p.Stanzas,
		ModelPath:      p.ModelPath,
		SourcetypePath: p.SourcetypePath,
		Families:       p.Families,
		IsPartial:      p.IsPartial,
		LocalStanzas:   p.LocalStanzas,
	}
}
