package entities

type Email struct {
	Email string `json:"email" xorm:"email" binding:"required,email"`
}
