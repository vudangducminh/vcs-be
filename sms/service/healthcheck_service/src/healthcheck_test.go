package src

import (
	"healthcheck_service/entities"
	"testing"
	"time"
)

type mockElasticQuery struct {
	servers []entities.BriefServerInfo
	updates []entities.ServerUptimeUpdate
}

func (m *mockElasticQuery) GetAllServer() []entities.BriefServerInfo {
	return m.servers
}
func (m *mockElasticQuery) BulkUpdateServerInfo(updates []entities.ServerUptimeUpdate) {
	m.updates = updates
}

type mockPinger struct {
	alive bool
}

func (m *mockPinger) PingServer(ip string) bool {
	return m.alive
}

func TestModifiedStartHealthCheck_ActiveServer(t *testing.T) {
	eq := &mockElasticQuery{
		servers: []entities.BriefServerInfo{
			{Id: "1", IPv4: "127.0.0.1", Uptime: []int{0}},
			{Id: "2", IPv4: "127.0.0.2", Uptime: []int{30}},
		},
	}
	pg := &mockPinger{alive: true}
	SetElasticQuery(eq)
	SetPinger(pg)

	// Run one iteration only for test
	done := make(chan bool)
	go func() {
		prevMinute := -1
		for i := 0; i < 1; i++ {
			if time.Now().Minute() != prevMinute {
				updateList := []entities.ServerUptimeUpdate{}
				ServerList := eq.GetAllServer()
				newUptimeStatus := time.Now().Minute()%20 == 0
				resultsChan := make(chan entities.ServerUptimeUpdate, len(ServerList))

				for _, server := range ServerList {
					go func(srv entities.BriefServerInfo) {
						uptime := srv.Uptime
						if newUptimeStatus {
							uptime = append(uptime, 0)
							if len(uptime) > 2016 {
								uptime = uptime[len(uptime)-2016:]
							}
						}
						isAlive := pg.PingServer(srv.IPv4)
						if isAlive {
							uptime[len(uptime)-1] += 60
						}
						status := "inactive"
						if isAlive {
							status = "active"
						}
						resultsChan <- entities.ServerUptimeUpdate{
							Id:     srv.Id,
							Uptime: uptime,
							Status: status,
						}
					}(server)
				}

				for i := 0; i < len(ServerList); i++ {
					update := <-resultsChan
					updateList = append(updateList, update)
				}
				close(resultsChan)
				eq.BulkUpdateServerInfo(updateList)
				prevMinute = time.Now().Minute()
			}
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()
	<-done

	if len(eq.updates) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(eq.updates))
	}
	if eq.updates[0].Status != "active" {
		t.Errorf("Expected status 'active', got %s", eq.updates[0].Status)
	}
}
