package healthcheckservice

import (
	"log"
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
	for {
		for _, server := range ServerList {
			// ping all server
			isAlive := PingServer(server.IPv4)
			if isAlive {
				log.Println("Server", server.ServerId, "is alive")
				// update data to elasticsearch
			} else {
				log.Println("Server", server.ServerId, "is not alive")
				// update data to elasticsearch
			}
		}
		time.Sleep(30 * time.Second)
	}

}
