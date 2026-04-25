package httpapi

// CreateRequest is the JSON body for POST /api/shorten.
type CreateRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

// CreateResponse is the JSON body returned for a created or reused short link.
type CreateResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
	URL      string `json:"url"`
}

// ErrorResponse is the JSON body returned for any non-2xx API response.
type ErrorResponse struct {
	Error string `json:"error"`
}
