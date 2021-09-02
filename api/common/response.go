package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type MetaData struct {
	Page       int64                 `json:"page"`
	PerPage    int64                 `json:"per_page"`
	PageCount  int64                 `json:"page_count"`
	TotalCount int64                 `json:"total_count"`
	Links      []map[string]string `json:"links"`
}

type ResponseDTO struct {
	Metadata *MetaData   `json:"_metadata"`
	Data     interface{} `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

func GenerateSuccessResponse(c echo.Context, data interface{}, metadata *MetaData, message string) error {
	if metadata != nil {
		return c.JSON(http.StatusOK, ResponseDTO{
			Status:   "success",
			Message:  message,
			Data:     data,
			Metadata: metadata,
		})
	}
	return c.JSON(http.StatusOK, ResponseDTO{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func GenerateErrorResponse(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusBadRequest, ResponseDTO{
		Status:  "error",
		Message: message,
		Data:    data,
	})
}
func GetPaginationMetadata(page, limit, totalRecords, totalPaginatedRecords int64) MetaData {
	metaData := MetaData{
		Page:       page,
		PerPage:    limit,
		TotalCount: totalRecords,
		PageCount:  totalPaginatedRecords,
	}
	return metaData
}
