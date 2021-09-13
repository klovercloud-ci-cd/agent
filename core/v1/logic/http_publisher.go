package logic

import (
	"bytes"
	"github.com/klovercloud-ci/core/v1/service"
	"io/ioutil"
	"log"
	"net/http"
)

type httpPublisherService struct {

}

func (h httpPublisherService) Post(url string, header map[string]string, body []byte) error{
	log.Println("Posting ...", url, string(body))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	for k,v:=range header{
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("[ERROR] Failed communicate agent:", err.Error())
		return err
	}else if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR] Failed communicate agent:", err.Error())
		} else {
			log.Println("[ERROR] Failed communicate agent::", string(body))
		}
	}
	return nil
}

func NewHttpPublisherService() service.HttpPublisher {
	return &httpPublisherService{
	}
}