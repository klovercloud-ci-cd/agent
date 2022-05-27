package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/agent/api/common"
	"github.com/klovercloud-ci-cd/agent/config"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/klovercloud-ci-cd/agent/enums"
	"log"
	"strconv"
)

type resourceService struct {
	K8s          service.K8s
	observerList []service.Observer
	httpClient   service.HttpClient
}

func (r resourceService) Pull() {
	url := config.ApiServiceUrl + "/process_life_cycle_events?count=" + config.PullSize + "&agent=" + config.AgentName
	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Authorization"] = "Bearer " + config.Token
	data, err := r.httpClient.Get(url, header)
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
	for _, each := range resources {
		err := r.Update(each)
		subject := v1.Subject{each.Step, "", each.Name, each.Namespace, each.ProcessId, nil, nil, nil}
		subject.EventData=make(map[string]interface{})
		subject.EventData["step"] = each.Name
		subject.EventData["process_id"] = each.ProcessId
		subject.EventData["company_id"] =  each.Pipeline.MetaData.CompanyId
		subject.EventData["claim"] = strconv.Itoa(each.Claim)
		subject.EventData["footmark"] = enums.POST_AGENT_JOB
		if err != nil {
			subject.Log = "Update Failed: " + err.Error()
			subject.EventData["log"] = subject.Log
			subject.EventData["status"] = enums.FAILED
			go r.notifyAll(subject)
		}else {
			subject.EventData["log"] = "Agent Job Completed"
			subject.EventData["status"] = enums.COMPLETED
			go r.notifyAll(subject)
		}
	}
}

func (r resourceService) Update(resource v1.Resource) error {
	listener := v1.Subject{Log: "Deploy Step Started", ProcessId: resource.ProcessId, Step: resource.Step}
	processEventData := make(map[string]interface{})
	processEventData["step"] = resource.Step
	processEventData["type"] = resource.Type
	processEventData["process_id"] = resource.ProcessId
	processEventData["company_id"] =  resource.Pipeline.MetaData.CompanyId
	processEventData["footmark"] = enums.INIT_AGNET_JOB
	processEventData["claim"] = strconv.Itoa(resource.Claim)
	listener.EventData = processEventData
	go r.notifyAll(listener)
	for _, each := range *resource.Descriptors {
		each.SetLabels(map[string]string{"company": resource.Pipeline.MetaData.CompanyId, "klovercloud_ci": "enabled","process_id":resource.ProcessId,"claim":strconv.Itoa(resource.Claim)})
		err:=r.K8s.Apply(each)
		if err!=nil{
			listener.Log=err.Error()
			go r.notifyAll(listener)
		}
	}
	if resource.Name == "" {
		subject := v1.Subject{resource.Step, "Updated Successfully", resource.Name, resource.Namespace, resource.ProcessId, map[string]interface{}{"footmark": enums.POST_AGENT_JOB, "log": "Updated Successfully", "reason": "n/a", "step":resource.Step,"process_id":resource.ProcessId,"company_id":resource.Pipeline.MetaData.CompanyId,"status": enums.SUCCESSFUL, "claim": strconv.Itoa(resource.Claim)}, nil, resource.Pipeline}
		go r.notifyAll(subject)
		return nil
	}
	subject := v1.Subject{resource.Step, "Updating resource", resource.Name, resource.Namespace, resource.ProcessId, map[string]interface{}{"footmark": enums.UPDATE_RESOURCE, "log": "Updating resource", "reason": "n/a","step":resource.Step,"process_id":resource.ProcessId,"company_id":resource.Pipeline.MetaData.CompanyId, "status": enums.SUCCESSFUL, "claim": strconv.Itoa(resource.Claim)}, nil, resource.Pipeline}
	go r.notifyAll(subject)
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
func (r resourceService) notifyAll(subject v1.Subject) {
	for _, observer := range r.observerList {
		go observer.Listen(subject)
	}
}

// NewResourceService returns resource type service.
func NewResourceService(k8s service.K8s, observerList []service.Observer, httpClient service.HttpClient) service.Resource {
	return &resourceService{
		K8s:          k8s,
		observerList: observerList,
		httpClient:   httpClient,
	}
}
