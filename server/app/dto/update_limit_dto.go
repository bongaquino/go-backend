package dto

type UpdateLimitDTO struct {
	BytesLimit int64 `json:"bytes_limit" binding:"required"`
}
