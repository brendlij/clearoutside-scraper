package api

type errorResponse struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type versionResponse struct {
	Version string `json:"version"`
}
