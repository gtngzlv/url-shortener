package models

type APIShortenRequest struct {
	URL string `json:"url"`
}

type APIShortenResponse struct {
	Result string `json:"result"`
}

type URLInfo struct {
	UUID          string `json:"uuid,omitempty"`
	UserID        string `json:"userID,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	IsDeleted     int    `json:"is_deleted,omitempty"`
}

type BatchRequest struct {
	Entities []URLInfo
}

type BatchResponse struct {
	Entities []URLInfo
}
