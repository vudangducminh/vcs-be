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
	// Try common server ports
	ports := []string{"22", "80", "443", "8080", "3389", "21", "25"}
	timeout := time.Second * 3

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
		if err == nil {
			conn.Close()
			return true // Server is reachable on at least one port
		}
	}
	return false // No ports were reachable
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

	var prevMinute int = -1

	for {
		if time.Now().Minute() != prevMinute {
			// Only proceed if we have servers to check
			updateList := []object.ServerUptimeUpdate{}
			log.Println("Attempting to refresh from Elasticsearch...")
			ServerList = elastic_query.GetAllServer()
			var newUptimeStatus bool = false
			if time.Now().Minute()%1 == 0 {
				newUptimeStatus = true
			}

			// Use a channel to collect results from goroutines
			resultsChan := make(chan object.ServerUptimeUpdate, len(ServerList))

			for _, server := range ServerList {
				// Ping all server
				go func(srv object.BriefServerInfo) {
					var uptime []int = srv.Uptime
					if newUptimeStatus {
						uptime = append(uptime, 0)
						// Save data for 1 week only
						if len(uptime) > 504 {
							uptime = uptime[len(uptime)-504:]
						}
					}
					isAlive := PingServer(srv.IPv4)
					if isAlive {
						uptime[len(uptime)-1] += 30
					}
					if uptime[0] > 0 {
						log.Println("IP", srv.IPv4, "uptime:", uptime)
					}

					// Send result to channel
					var status string
					if isAlive {
						status = "active"
					} else {
						status = "inactive"
					}
					resultsChan <- object.ServerUptimeUpdate{
						ServerId: srv.ServerId,
						Uptime:   uptime,
						Status:   status,
					}
				}(server)
			}

			// Collect all results

			// Explanation of how this works:

			// 1) Channel buffer contains (from goroutines):
			// resultsChan: [struct1, struct2, struct3]

			// 2) First receive:
			// update := <-resultsChan
			// Now: update = struct1, channel: [struct2, struct3]

			// 3) Second receive (next loop iteration):
			// update := <-resultsChan
			// Now: update = struct2, channel: [struct3]

			// 4) Third receive:
			// update := <-resultsChan
			// Now: update = struct3, channel: []
			for i := 0; i < len(ServerList); i++ {
				update := <-resultsChan
				updateList = append(updateList, update)
			}
			close(resultsChan) // Close the channel when done

			log.Println("Update list length:", len(updateList))
			elastic_query.BulkUpdateServerInfo(updateList)
			prevMinute = time.Now().Minute()
		}
		time.Sleep(10 * time.Second)
	}
}
