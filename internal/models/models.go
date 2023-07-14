package models

type APIShortenRequest struct {
	URL string `json:"url"`
}

type APIShortenResponse struct {
	Result string `json:"result"`
}

type BatchEntity struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

type BatchRequest struct {
	Entities []BatchEntity
}

type BatchResponse struct {
	Entities []BatchEntity
}
