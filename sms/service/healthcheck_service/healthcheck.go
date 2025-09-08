package healthcheckservice

import (
	"net"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"
)

var ServerList []object.BriefServerInfo

func PingServer(ip string) bool {
	timeout := time.Second * 3
	conn, err := net.DialTimeout("ip4:icmp", ip, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func HealthCheck() {
	if len(ServerList) == 0 {
		ServerList = elastic_query.GetAllServer()
	}
	for _, server := range ServerList {
		// ping all server
		isAlive := PingServer(server.IPv4)
		if isAlive {
			// update data to elasticsearch
		} else {
			// update data to elasticsearch
		}
	}
}
