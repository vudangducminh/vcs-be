package healthcheck_service

import (
	"context"
	"log"
	"net"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"
)

var ServerList []object.BriefServerInfo

func PingServer(ip string) bool {
	timeout := time.Second * 3

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: timeout,
			}
			return d.DialContext(ctx, network, address)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := resolver.LookupAddr(ctx, ip)
	return err == nil
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
			go func(ip string) {
				isAlive := PingServer(ip)
				log.Println("IP", ip, "alive status:", isAlive)
			}(server.IPv4)
		}
		// log.Println("Pinging 196.56.217.125 inside container")
		// log.Println("Status: ", PingServer("196.56.217.125"))
		time.Sleep(30 * time.Second)
	}

}
