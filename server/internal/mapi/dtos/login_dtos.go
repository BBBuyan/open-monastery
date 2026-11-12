package dtos

type LoginRequest struct {
	UserID   int64  `json:"user_id" binding:"required,gt=0"`
	Password string `json:"password" binding:"required"`
}
