package src

import (
	"report_service/entities"
	"testing"

	"github.com/xuri/excelize/v2"
)

type mockDeps struct {
	emailStatus     int
	serverStatus    int
	sendEmailStatus int
	emails          []entities.Email
	servers         []entities.Server
	avgUptime       float64
}

func (m *mockDeps) GetAllEmails() ([]entities.Email, int) {
	return m.emails, m.emailStatus
}
func (m *mockDeps) GetServerUptimeInRange(start, end int64, order, filter, substr string) ([]entities.Server, int, float64) {
	return m.servers, m.serverStatus, m.avgUptime
}
func (m *mockDeps) GetTotalActiveServersCount(order, filter string) int      { return 1 }
func (m *mockDeps) GetTotalInactiveServersCount(order, filter string) int    { return 1 }
func (m *mockDeps) GetTotalMaintenanceServersCount(order, filter string) int { return 1 }
func (m *mockDeps) SendEmail(f *excelize.File, to, subject, body string) int {
	return m.sendEmailStatus
}

func TestModifiedDailyReporter_Success(t *testing.T) {
	deps := &mockDeps{
		emailStatus:     200,
		serverStatus:    200,
		sendEmailStatus: 200,
		emails:          []entities.Email{{Email: "test@example.com"}},
		servers:         []entities.Server{{Id: "1", ServerName: "srv", Status: "active", IPv4: "127.0.0.1", Uptime: []int{100}, CreatedTime: 1, LastUpdatedTime: 2}},
		avgUptime:       99.9,
	}
	code, err := ModifiedDailyReporter(deps, 100000)
	if code != 200 || err != nil {
		t.Errorf("Expected 200, got %d, err: %v", code, err)
	}
}

func TestModifiedDailyReporter_EmailFail(t *testing.T) {
	deps := &mockDeps{emailStatus: 500}
	code, err := ModifiedDailyReporter(deps, 100000)
	if code != 500 || err == nil {
		t.Errorf("Expected 500 and error, got %d, err: %v", code, err)
	}
}

func TestModifiedDailyReporter_NoEmails(t *testing.T) {
	deps := &mockDeps{emailStatus: 200, emails: []entities.Email{}}
	code, err := ModifiedDailyReporter(deps, 100000)
	if code != 204 || err == nil {
		t.Errorf("Expected 204 and error, got %d, err: %v", code, err)
	}
}

func TestModifiedDailyReporter_ServerFail(t *testing.T) {
	deps := &mockDeps{emailStatus: 200, emails: []entities.Email{{Email: "test@example.com"}}, serverStatus: 500}
	code, err := ModifiedDailyReporter(deps, 100000)
	if code != 500 || err == nil {
		t.Errorf("Expected 500 and error, got %d, err: %v", code, err)
	}
}

func TestModifiedDailyReporter_SendEmailFail(t *testing.T) {
	deps := &mockDeps{
		emailStatus:     200,
		serverStatus:    200,
		sendEmailStatus: 500,
		emails:          []entities.Email{{Email: "test@example.com"}},
		servers:         []entities.Server{{Id: "1", ServerName: "srv", Status: "active", IPv4: "127.0.0.1", Uptime: []int{100}, CreatedTime: 1, LastUpdatedTime: 2}},
		avgUptime:       99.9,
	}
	code, err := ModifiedDailyReporter(deps, 100000)
	if code != 500 || err == nil {
		t.Errorf("Expected 500 and error, got %d, err: %v", code, err)
	}
}
