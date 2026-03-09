package response

// HTTPError represents a standard error response
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
	Details string `json:"details,omitempty" example:"missing required field 'guestId'"`
}
