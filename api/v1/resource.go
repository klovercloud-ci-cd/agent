package v1

import (
	"github.com/klovercloud-ci/api/common"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/api"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/labstack/echo/v4"
	"log"
)

type resourceApi struct {
	resourceService service.Resource
}

func (r resourceApi) Update(context echo.Context) error {
	data:=v1.Resource{}
	err := context.Bind(&data)
	if  err != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,nil,err.Error())
	}
	error :=r.resourceService.Update(data)

	if error != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,err.Error(),"Failed to trigger deploy!")
	}
	return common.GenerateSuccessResponse(context,"",nil,"Pipeline successfully triggered!")
}

func NewResourceApi(resourceService service.Resource) api.Resource {
	return &resourceApi{
		resourceService: resourceService,
	}
}
