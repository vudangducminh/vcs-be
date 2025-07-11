package object

type Server struct {
	ServerId        string `xorm:"'server_id' pk"`
	ServerName      string `xorm:"'server_name'"`
	Status          string `xorm:"'status'"`            // e.g., "active", "inactive", "maintenance"
	Uptime          int    `xorm:"'uptime'"`            // e.g., "3666" for 1 hour 1 minute and 6 seconds
	CreatedTime     string `xorm:"'created_time'"`      // ISO 8601 format
	LastUpdatedTime string `xorm:"'last_updated_time'"` // ISO 8601 format
	IPv4            string `xorm:"'ipv4'"`
}
