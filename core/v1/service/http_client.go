package service

// HttpClient HttpClient operations.
type HttpClient interface {
	Post(url string, header map[string]string, body []byte) error
	Get(url string, header map[string]string) ([]byte, error)
}
