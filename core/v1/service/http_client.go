package service

type HttpClient interface {
	Post(url string,header map[string]string,body []byte) error
	Get(url string, header map[string]string) (error, []byte)
}
