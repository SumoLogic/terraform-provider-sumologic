package sumologic

type RumSource struct {
	HTTPSource
	RumThirdPartyRef RumThirdPartyRef `json:"thirdPartyRef"`
}

type RumThirdPartyRef struct {
	Resources []RumThirdPartyResource `json:"resources"`
}

type RumThirdPartyResource struct {
	ServiceType    string                `json:"serviceType"`
	Authentication PollingAuthentication `json:"authentication,omitempty"`
	Path           RumSourcePath         `json:"path,omitempty"`
}

type RumSourcePath struct {
	Type                         string                 `json:"type"`
	ApplicationName              string                 `json:"applicationName"`
	ServiceName                  string                 `json:"serviceName"`
	DeploymentEnvironment        string                 `json:"deploymentEnvironment"`
	SamplingRate                 float32                `json:"samplingRate"`
	IgnoreUrls                   []string               `json:ignoreUrls"`
	CustomTags                   map[string]interface{} `json:"customTags"`
	PropagateTraceHeaderCorsUrls []string               `json:"propagateTraceHeaderCorsUrls,omitempty"`
	SelectedCountry              string                 `json:"selectedCountry,omitempty"`
}
