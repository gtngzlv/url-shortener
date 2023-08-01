package models

type APIShortenRequest struct {
	URL string `json:"url"`
}

type APIShortenResponse struct {
	Result string `json:"result"`
}

type BatchEntity struct {
	UserID        string `json:"userID,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

type BatchRequest struct {
	Entities []BatchEntity
}

type BatchResponse struct {
	Entities []BatchEntity
}
