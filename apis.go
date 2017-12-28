package egoscale

/*

APIs

All the available APIs on the server and provided by the API Discovery plugin

	cs := egoscale.NewClient("https://api.exoscale.ch/compute", "EXO...", "...")

	resp := new(egoscale.ListAPIsResponse)
	err := cs.Request(&egoscale.ListAPIsRequest{}, resp)
	if err != nil {
		panic(err)
	}

	for _, api := range resp.API {
		fmt.Println("%s %s", api.Name, api.Description)
	}
	// Output:
	// listNetworks Lists all available networks
	// ...

*/

// API represents an API service
type API struct {
	Description string         `json:"description"`
	IsAsync     bool           `json:"isasync"`
	Name        string         `json:"name"`
	Related     string         `json:"related"` // comma separated
	Since       string         `json:"since"`
	Type        string         `json:"type"`
	Params      []*APIParam    `json:"params"`
	Response    []*APIResponse `json:"responses"`
}

// APIParam represents an API parameter field
type APIParam struct {
	Description string `json:"description"`
	Length      int64  `json:"length"`
	Name        string `json:"name"`
	Related     string `json:"related"` // comma separated
	Since       string `json:"since"`
	Type        string `json:"type"`
}

// APIResponse represents an API response field
type APIResponse struct {
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Response    []*APIResponse `json:"response"`
	Type        string         `json:"type"`
}

// ListAPIsRequest represents a query to list the api
type ListAPIsRequest struct {
	Name string `json:"name,omitempty"`
}

func (req *ListAPIsRequest) name() string {
	return "listAPIs"
}

func (req *ListAPIsRequest) response() interface{} {
	return new(ListAPIsResponse)
}

// ListAPIsResponse represents a list of API
type ListAPIsResponse struct {
	Count int    `json:"count"`
	API   []*API `json:"api"`
}
