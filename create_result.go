package goelasticsearch

// CreateResult represents a create result.
type CreateResult struct {
	Created bool   `json:"created"`
	Error   string `json:"error"`
	ID      string `json:"_id"`
	Index   string `json:"_index"`
	Status  int16  `json:"status"`
	Type    string `json:"_type"`
	Version int64  `json:"_version"`
}
