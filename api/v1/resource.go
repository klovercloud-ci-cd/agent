package v1

import (
	"github.com/klovercloud-ci-cd/agent/api/common"
	v1 "github.com/klovercloud-ci-cd/agent/core/v1"
	"github.com/klovercloud-ci-cd/agent/core/v1/api"
	"github.com/klovercloud-ci-cd/agent/core/v1/service"
	"github.com/labstack/echo/v4"
	"log"
)

type resourceApi struct {
	resourceService service.Resource
}

// Apply... Apply resources
// @Summary Apply resources
// @Description Apply resources
// @Tags Resource
// @Produce json
// @Param data body v1.Resource true "Resource Data"
// @Success 200 {object} common.ResponseDTO
// @Router /api/v1/resources/{processId} [POST]
func (r resourceApi) Update(context echo.Context) error {
	data := v1.Resource{}
	err := context.Bind(&data)
	if err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, err.Error())
	}
	error := r.resourceService.Update(data)

	if error != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, err.Error(), "Failed to trigger deploy!")
	}
	return common.GenerateSuccessResponse(context, "", nil, "Pipeline successfully triggered!")
}

// NewResourceApi returns Resource type api
func NewResourceApi(resourceService service.Resource) api.Resource {
	return &resourceApi{
		resourceService: resourceService,
	}
}
