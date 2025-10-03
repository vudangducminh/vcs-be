package object

type ServerUptimeUpdate struct {
	Id     string `json:"_id" binding:"required"`
	Uptime []int  `json:"uptime" binding:"required"`
	Status string `json:"status,omitempty"`
}
