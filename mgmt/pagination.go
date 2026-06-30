package mgmt

// Pagination is the standard SentinelOne pagination envelope.
type Pagination struct {
	TotalItems int    `json:"totalItems"`
	NextCursor string `json:"nextCursor"`
}
