package healthcheck_service

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
	// Wait a bit for services to fully start
	time.Sleep(10 * time.Second)

	if len(ServerList) == 0 {
		log.Println("Attempting to fetch server list from Elasticsearch...")
		ServerList = elastic_query.GetAllServer()
		if len(ServerList) == 0 {
			log.Println("No servers found or Elasticsearch not ready, will retry later")
		}
	}
	for {
		// Only proceed if we have servers to check
		if len(ServerList) == 0 {
			log.Println("Server list is empty, attempting to refresh from Elasticsearch...")
			ServerList = elastic_query.GetAllServer()
		}

		for _, server := range ServerList {
			// ping all server
			isAlive := PingServer(server.IPv4)
			if isAlive {
				// log.Println("IP", server.IPv4, "is alive")
				// update data to elasticsearch
			} else {
				// log.Println("IP", server.IPv4, "is not alive")
				// update data to elasticsearch
			}
		}
		time.Sleep(30 * time.Second)
	}

}
