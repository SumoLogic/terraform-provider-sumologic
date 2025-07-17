package sumologic

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (s *Client) GetCSELogMappingVendorsAndProducts(product, vendor string) (*CSELogMappingVendorAndProduct, error) {

	data, err := s.Get(fmt.Sprintf("sec/v1/vendors-and-products?product=%s&vendor=%s", url.QueryEscape(product), url.QueryEscape(vendor)))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSELogMappingVendorAndProductResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	productAndVendorSlice := response.CSELogMappingVendorAndProduct
	var productAndVendor CSELogMappingVendorAndProduct

	if len(productAndVendorSlice) == 1 {
		productAndVendor = productAndVendorSlice[0]
	}

	return &productAndVendor, nil
}

type CSELogMappingVendorAndProductResponse struct {
	CSELogMappingVendorAndProduct []CSELogMappingVendorAndProduct `json:"data"`
}

type CSELogMappingVendorAndProduct struct {
	GUID    string `json:"guid"`
	Product string `json:"product"`
	Vendor  string `json:"vendor"`
}
