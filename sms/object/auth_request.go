package object

type AuthRequest struct {
	JWT string `json:"jwt" binding:"required"`
}
