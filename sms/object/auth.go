package object

type AuthRequest struct {
	JWT string `json:"jwt" binding:"required"`
}

type AuthResponse struct {
	Message string `json:"message" example:"Authentication successfully"`
}
