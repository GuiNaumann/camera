package entities

// SuccessfulRequest response for simple successful requests
type SuccessfulRequest struct {
	// Success flag indicating the success
	Success bool `json:"success"`
}

func NewSuccessfulRequest() SuccessfulRequest {
	return SuccessfulRequest{
		Success: true,
	}
}
