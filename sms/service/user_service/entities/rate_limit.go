package entities

type RateLimitExceededResponse struct {
	Error string `json:"error" example:"Rate limit exceeded. Please try again later."`
}
