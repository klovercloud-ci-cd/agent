package logic

import (
	"bytes"
	"errors"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type httpClientService struct {
}

func (h httpClientService) Post(url string, header map[string]string, body []byte) error {
	log.Println("Posting ...", url, string(body))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[ERROR] Failed communicate agent:", err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR] Failed communicate agent:", err.Error())
		} else {
			log.Println("[ERROR] Failed communicate agent::", string(body))
		}
	}
	return nil
}

func (h httpClientService) Get(url string, header map[string]string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		jsonDataFromHttp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return jsonDataFromHttp, nil
	}
	return nil, errors.New("Status: " + res.Status + ", code: " + strconv.Itoa(res.StatusCode))
}

// NewHttpClientService returns HttpClient type service
func NewHttpClientService() service.HttpClient {
	return &httpClientService{}
}