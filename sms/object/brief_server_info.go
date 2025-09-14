package object

type BriefServerInfo struct {
	ServerId string `json:"server_id"`
	IPv4     string `json:"ipv4"`
	Uptime   []int  `json:"uptime"`
}
