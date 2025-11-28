package src

import (
	"fmt"
	"log"
	"net/http"
	"report_service/entities"
	elastic_query "report_service/infrastructure/elasticsearch/query"
	postgresql_query "report_service/infrastructure/postgresql/query"
	"report_service/src/template"
	"time"

	"github.com/xuri/excelize/v2"
)

func DailyReporter() {
	time.Sleep(10 * time.Second) // Initial delay to allow other services to start

	for {
		var sec int64 = 86400
		if time.Now().Unix()%sec >= sec-20 {
			log.Println("Starting daily report email process...")
			emailList, status := postgresql_query.GetAllEmails()
			if status != http.StatusOK {
				log.Println("Failed to retrieve email list from PostgreSQL")
			} else {
				if len(emailList) == 0 {
					log.Println("No email addresses found in PostgreSQL")
					time.Sleep(20 * time.Second)
					continue
				}
				currentTimeInSecond := time.Now().Unix()
				serverDataList, status, averageUptimePercentage := elastic_query.GetServerUptimeInRange(currentTimeInSecond-86400, currentTimeInSecond, "", "", "")
				totalActiveServer := elastic_query.GetTotalActiveServersCount("", "")
				totalInactiveServer := elastic_query.GetTotalInactiveServersCount("", "")
				totalMaintenanceServer := elastic_query.GetTotalMaintenanceServersCount("", "")
				totalServer := totalActiveServer + totalInactiveServer + totalMaintenanceServer
				log.Printf("Total servers: %d, Active: %d, Inactive: %d, Maintenance: %d, Average Uptime: %.2f%%\n", totalServer, totalActiveServer, totalInactiveServer, totalMaintenanceServer, averageUptimePercentage)
				if status != http.StatusOK {
					log.Println("Failed to retrieve server details from Elasticsearch")
					time.Sleep(20 * time.Second)
					continue
				}
				f := excelize.NewFile()
				sheet := "Sheet1"
				f.SetSheetName("Sheet1", sheet)

				// Write header
				headers := []string{"Index", "Server ID", "Server Name", "Status", "IPv4", "Uptime", "Created Time", "Last Updated Time"}
				for i, h := range headers {
					cell, _ := excelize.CoordinatesToCellName(i+1, 1)
					f.SetCellValue(sheet, cell, h)
				}

				// Write data

				for rowIdx, server := range serverDataList {
					// Convert timestamps to readable format
					createdTimeStr := time.Unix(server.CreatedTime, 0).Format("2006-01-02 15:04:05")
					lastUpdatedTimeStr := time.Unix(server.LastUpdatedTime, 0).Format("2006-01-02 15:04:05")
					serverUptime := time.Unix(int64(server.Uptime[0]), 0).Format("15:04:05")
					values := []interface{}{
						rowIdx + 1,
						server.Id,
						server.ServerName,
						server.Status,
						server.IPv4,
						serverUptime,
						createdTimeStr,
						lastUpdatedTimeStr,
					}
					for colIdx, value := range values {
						cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
						f.SetCellValue(sheet, cell, value)
					}
				}
				emailBody := "Here is your requested server report." + "\n"
				emailBody += "Total servers in the system: " + fmt.Sprintf("%d", totalServer) + "\n"
				emailBody += "Number of active servers: " + fmt.Sprintf("%d", totalActiveServer) + "\n"
				emailBody += "Number of inactive servers: " + fmt.Sprintf("%d", totalInactiveServer) + "\n"
				emailBody += "Number of maintenance servers: " + fmt.Sprintf("%d", totalMaintenanceServer) + "\n"
				emailBody += "Average uptime percentage across all servers: " + fmt.Sprintf("%.2f", averageUptimePercentage) + "%" + "\n"

				for _, email := range emailList {
					// Send email with the Excel file as attachment
					status = template.SendEmail(f, email.Email, "Server Report", emailBody)
					if status != http.StatusOK {
						log.Printf("Failed to send daily report to %s\n", email.Email)
					} else {
						log.Printf("Daily report sent to %s\n", email.Email)
					}
				}
			}
		}
		time.Sleep(20 * time.Second) // Check every 20 seconds
	}
}

type DailyReportDeps interface {
	GetAllEmails() ([]entities.Email, int)
	GetServerUptimeInRange(start, end int64, order, filter, substr string) ([]entities.Server, int, float64)
	GetTotalActiveServersCount(order, filter string) int
	GetTotalInactiveServersCount(order, filter string) int
	GetTotalMaintenanceServersCount(order, filter string) int
	SendEmail(f *excelize.File, to, subject, body string) int
}

func ModifiedDailyReporter(deps DailyReportDeps, now int64) (int, error) {
	emailList, status := deps.GetAllEmails()
	if status != 200 {
		return status, fmt.Errorf("failed to get emails")
	}
	if len(emailList) == 0 {
		return 204, fmt.Errorf("no emails")
	}
	serverDataList, status, avgUptime := deps.GetServerUptimeInRange(now-86400, now, "", "", "")
	if status != 200 {
		return status, fmt.Errorf("failed to get server data")
	}
	totalActive := deps.GetTotalActiveServersCount("", "")
	totalInactive := deps.GetTotalInactiveServersCount("", "")
	totalMaintenance := deps.GetTotalMaintenanceServersCount("", "")
	totalServer := totalActive + totalInactive + totalMaintenance

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"Index", "Server ID", "Server Name", "Status", "IPv4", "Uptime", "Created Time", "Last Updated Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for rowIdx, server := range serverDataList {
		values := []interface{}{
			rowIdx + 1,
			server.Id,
			server.ServerName,
			server.Status,
			server.IPv4,
			server.Uptime[0],
			server.CreatedTime,
			server.LastUpdatedTime,
		}
		for colIdx, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}
	emailBody := fmt.Sprintf(
		"Here is your requested server report.\nTotal servers in the system: %d\nNumber of active servers: %d\nNumber of inactive servers: %d\nNumber of maintenance servers: %d\nAverage uptime percentage across all servers: %.2f%%\n",
		totalServer, totalActive, totalInactive, totalMaintenance, avgUptime,
	)
	for _, email := range emailList {
		status = deps.SendEmail(f, email.Email, "Server Report", emailBody)
		if status != 200 {
			return status, fmt.Errorf("failed to send email to %s", email.Email)
		}
	}
	return 200, nil
}
