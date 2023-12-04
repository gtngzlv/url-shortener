package models

// APIShortenRequest model for /api/shorten request
type APIShortenRequest struct {
	URL string `json:"url"`
}

// APIShortenResponse model for /api/shorten response
type APIShortenResponse struct {
	Result string `json:"result"`
}

// URLInfo model for url info
type URLInfo struct {
	UUID          string `json:"uuid,omitempty"`
	UserID        string `json:"userID,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	IsDeleted     int    `json:"is_deleted,omitempty"`
}

// BatchRequest model for batch request
type BatchRequest struct {
	Entities []URLInfo
}

// BatchResponse model for batch response
type BatchResponse struct {
	Entities []URLInfo
}
