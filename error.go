package pi

import "strconv"

// HTTPError represents a HTTP Error.
type HTTPError struct {
	StatusCode int64 `json:"statusCode"`
}

func (h HTTPError) Error() string {
	return "error: " + strconv.FormatInt(h.StatusCode, 10)
}
