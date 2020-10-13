package option

// Response is the return value of the service action.
type Response struct {
	// Key represents an options header key.
	Key string `json:"key"`
	// Value represents an options header value.
	Value string `json:"value"`
}

// DefaultResponse provides a default response by best effort.
func DefaultResponse() []*Response {
	return []*Response{}
}
