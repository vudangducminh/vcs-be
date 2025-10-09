package src

import (
	"healthcheck_service/entities"
	elastic_query "healthcheck_service/infrastructure/elasticsearch/query"
	"log"
	"os/exec"
	"runtime"
	"time"
)

var ServerList []entities.BriefServerInfo

// Ping an IP address by attempting to open a TCP connection on common ports
func PingServer(ip string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "3", ip)
	}

	err := cmd.Run()
	return err == nil
}

func StartHealthCheck() {
	log.Println("Starting Health Check logic...")
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
			updateList := []entities.ServerUptimeUpdate{}
			log.Println("Attempting to refresh from Elasticsearch...")
			ServerList = elastic_query.GetAllServer()
			var newUptimeStatus bool = false
			if time.Now().Minute()%20 == 0 {
				newUptimeStatus = true
			}

			// Use a channel to collect results from goroutines
			resultsChan := make(chan entities.ServerUptimeUpdate, len(ServerList))

			for _, server := range ServerList {
				// Ping all server
				go func(srv entities.BriefServerInfo) {
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
						uptime[len(uptime)-1] += 60
						// log.Printf("IP %s is alive\n", srv.IPv4)
					}

					// Send result to channel
					var status string
					if isAlive {
						status = "active"
					} else {
						status = "inactive"
					}
					resultsChan <- entities.ServerUptimeUpdate{
						Id:     srv.Id,
						Uptime: uptime,
						Status: status,
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
