package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci/api/common"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"log"
)

type resourceService struct {
	K8s          service.K8s
	observerList []service.Observer
	httpClient service.HttpClient
}

func (r resourceService) Pull() {
	url := config.EventStoreUrl+"/process_life_cycle_events?count="+config.PullSize+"&agent="+config.AgentName
	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["token"]=config.Token
	err, data := r.httpClient.Get(url,header)
	if err != nil {
		// send to observer
		log.Println(err.Error())
		return
	}
	response := common.ResponseDTO{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return
	}
	b, err := json.Marshal(response.Data)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return
	}
	resources := []v1.Resource{}
	err = json.Unmarshal(b, &resources)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return
	}
	for _,each:=range resources{
		r.Update(each)
	}
}

func (r resourceService) Update(resource v1.Resource) error {
	for _,each:= range *resource.Descriptors{
		r.K8s.Apply(each)
	}
	if resource.Type == enums.DEPLOYMENT {
		return r.K8s.UpdateDeployment(resource)
	} else if resource.Type == enums.POD {
		return r.K8s.UpdatePod(resource)
	} else if resource.Type == enums.STATEFULSET {
		return r.K8s.UpdateStatefulSet(resource)
	} else if resource.Type == enums.DAEMONSET {
		return r.K8s.UpdateDaemonSet(resource)
	}
	return nil
}
func (r resourceService)notifyAll(subject v1.Subject){
	for _, observer := range r.observerList {
		go observer.Listen(subject)
	}
}

func NewResourceService(k8s service.K8s, observerList []service.Observer,httpClient service.HttpClient) service.Resource {
	return &resourceService{
		K8s:          k8s,
		observerList: observerList,
		httpClient: httpClient,
	}
}
