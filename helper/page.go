package helper

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"page_size"`
	Count     int64 `json:"count"`
	TotalPage int64 `json:"total_page"`
	Offset    int64 `json:"-"`
}

func CreatePageParams(Page, PageSize int) *PageData {
	params := PageData{}
	if Page <= 0 {
		params.Page = 1
	}
	if PageSize <= 0 {
		params.PageSize = 10
	}
	params.Offset = (params.Page - 1) * params.PageSize
	return &params
}

func GetPageParams(c *gin.Context) *PageData {
	var params PageData
	params.Page, _ = strconv.ParseInt(c.Query("p"), 10, 64)
	params.PageSize, _ = strconv.ParseInt(c.Query("l"), 10, 64)

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	params.Offset = (params.Page - 1) * params.PageSize
	return &params
}

func SetPageData(Page, PageSize, Count int64) *PageData {
	TotalPage := math.Ceil(float64(Count) / float64(PageSize))
	if TotalPage < 1 {
		TotalPage = 1
	}
	return &PageData{
		Page:      Page,
		PageSize:  PageSize,
		Count:     Count,
		TotalPage: int64(TotalPage),
	}
}
