package httpx

type PaginationQuery struct {
	Limit  int32 `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset,default=0" binding:"omitempty,min=0"`
}

type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	TotalCount int64 `json:"totalCount"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}
