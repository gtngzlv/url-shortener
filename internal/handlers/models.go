package handlers

type APIShortenRequest struct {
	URL string `json:"url"`
}

type APIShortenResponse struct {
	Result string `json:"result"`
}
