package object

type ServerUptimeUpdate struct {
	ServerId string `json:"server_id" binding:"required"`
	Uptime   []int  `json:"uptime" binding:"required"`
	Status   string `json:"status,omitempty"`
}
