package object

type ViewServerResponse struct {
	ServerId        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`
	Uptime          string `json:"uptime"`
	CreatedTime     string `json:"created_time"`
	LastUpdatedTime string `json:"last_updated_time"`
	IPv4            string `json:"ipv4"`
}
