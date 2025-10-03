package object

type BriefServerInfo struct {
	Id     string `json:"_id"`
	IPv4   string `json:"ipv4"`
	Uptime []int  `json:"uptime"`
}
