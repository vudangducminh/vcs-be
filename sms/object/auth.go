package object

type AuthRequest struct {
	JWT string `json:"jwt" binding:"required"`
}

type AuthSuccessResponse struct {
	Message string `json:"message" example:"Authentication successfully"`
}

type AuthErrorResponse struct {
	Error string `json:"error" example:"Authentication failed"`
}
